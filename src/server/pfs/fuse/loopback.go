// Copyright 2019 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"context"
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"github.com/pachyderm/pachyderm/src/client"
	"github.com/pachyderm/pachyderm/src/client/pfs"
	"github.com/pachyderm/pachyderm/src/server/pkg/errutil"
)

type fileState int32

const (
	none    fileState = iota // we don't know about this file
	meta                     // we have meta information (but not content for this file)
	full                     // we have full content for this file
	dirty                    // we have full content for this file and the user has written to it
	deleted                  // this file has been deleted (and the local version reflects this)
)

type loopbackRoot struct {
	loopbackNode

	rootPath string
	rootDev  uint64

	c *client.APIClient

	branches map[string]string
	commits  map[string]string
	files    map[string]fileState
	mu       sync.Mutex
}

type loopbackNode struct {
	fs.Inode
}

var _ = (fs.NodeStatfser)((*loopbackNode)(nil))
var _ = (fs.NodeStatfser)((*loopbackNode)(nil))
var _ = (fs.NodeGetattrer)((*loopbackNode)(nil))
var _ = (fs.NodeGetxattrer)((*loopbackNode)(nil))
var _ = (fs.NodeSetxattrer)((*loopbackNode)(nil))
var _ = (fs.NodeRemovexattrer)((*loopbackNode)(nil))
var _ = (fs.NodeListxattrer)((*loopbackNode)(nil))
var _ = (fs.NodeReadlinker)((*loopbackNode)(nil))
var _ = (fs.NodeOpener)((*loopbackNode)(nil))
var _ = (fs.NodeCopyFileRanger)((*loopbackNode)(nil))
var _ = (fs.NodeLookuper)((*loopbackNode)(nil))
var _ = (fs.NodeOpendirer)((*loopbackNode)(nil))
var _ = (fs.NodeReaddirer)((*loopbackNode)(nil))
var _ = (fs.NodeMkdirer)((*loopbackNode)(nil))
var _ = (fs.NodeMknoder)((*loopbackNode)(nil))
var _ = (fs.NodeLinker)((*loopbackNode)(nil))
var _ = (fs.NodeSymlinker)((*loopbackNode)(nil))
var _ = (fs.NodeUnlinker)((*loopbackNode)(nil))
var _ = (fs.NodeRmdirer)((*loopbackNode)(nil))
var _ = (fs.NodeRenamer)((*loopbackNode)(nil))

func (n *loopbackNode) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(n.path(), &s)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStatfsT(&s)
	return fs.OK
}

func (n *loopbackRoot) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {

	st := syscall.Stat_t{}
	err := syscall.Stat(n.rootPath, &st)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStat(&st)
	return fs.OK
}

func (n *loopbackNode) root() *loopbackRoot {
	return n.Root().Operations().(*loopbackRoot)
}

func (n *loopbackNode) c() *client.APIClient {
	return n.root().c
}

func (n *loopbackNode) path() string {
	path := n.Path(nil)
	return filepath.Join(n.root().rootPath, path)
}

func (n *loopbackNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	if err := n.download(p, meta); err != nil {
		return nil, toErrno(err)
	}

	st := syscall.Stat_t{}
	err := syscall.Lstat(p, &st)
	if err != nil {
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))
	return ch, 0
}

func (n *loopbackNode) Mknod(ctx context.Context, name string, mode, rdev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	err := syscall.Mknod(p, mode, int(rdev))
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Rmdir(p)
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	return ch, 0
}

func (n *loopbackNode) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	if err := n.download(p, meta); err != nil {
		errno := toErrno(err)
		if errno != syscall.ENOENT {
			return nil, errno
		}
	}
	err := os.Mkdir(p, os.FileMode(mode))
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Rmdir(p)
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	return ch, 0
}

func (n *loopbackNode) Rmdir(ctx context.Context, name string) syscall.Errno {
	p := filepath.Join(n.path(), name)
	if err := n.download(p, meta); err != nil {
		return toErrno(err)
	}
	err := syscall.Rmdir(p)
	return fs.ToErrno(err)
}

func (n *loopbackNode) Unlink(ctx context.Context, name string) (errno syscall.Errno) {
	p := filepath.Join(n.path(), name)
	if err := n.download(p, meta); err != nil {
		return toErrno(err)
	}
	defer func() {
		if errno == 0 {
			n.root().mu.Lock()
			n.root().files[n.trimPath(p)] = deleted
			n.root().mu.Unlock()
		}
	}()
	err := syscall.Unlink(p)
	return fs.ToErrno(err)
}

func toLoopbackNode(op fs.InodeEmbedder) *loopbackNode {
	if r, ok := op.(*loopbackRoot); ok {
		return &r.loopbackNode
	}
	return op.(*loopbackNode)
}

func (n *loopbackNode) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	newParentLoopback := toLoopbackNode(newParent)
	if flags&fs.RENAME_EXCHANGE != 0 {
		return n.renameExchange(name, newParentLoopback, newName)
	}

	p1 := filepath.Join(n.path(), name)

	p2 := filepath.Join(newParentLoopback.path(), newName)
	err := os.Rename(p1, p2)
	return fs.ToErrno(err)
}

func (r *loopbackRoot) idFromStat(st *syscall.Stat_t) fs.StableAttr {
	// We compose an inode number by the underlying inode, and
	// mixing in the device number. In traditional filesystems,
	// the inode numbers are small. The device numbers are also
	// small (typically 16 bit). Finally, we mask out the root
	// device number of the root, so a loopback FS that does not
	// encompass multiple mounts will reflect the inode numbers of
	// the underlying filesystem
	swapped := (uint64(st.Dev) << 32) | (uint64(st.Dev) >> 32)
	swappedRootDev := (r.rootDev << 32) | (r.rootDev >> 32)
	return fs.StableAttr{
		Mode: uint32(st.Mode),
		Gen:  1,
		// This should work well for traditional backing FSes,
		// not so much for other go-fuse FS-es
		Ino: (swapped ^ swappedRootDev) ^ st.Ino,
	}
}

var _ = (fs.NodeCreater)((*loopbackNode)(nil))

func (n *loopbackNode) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (inode *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	p := filepath.Join(n.path(), name)
	if err := n.download(p, full); err != nil {
		errno := toErrno(err)
		if errno != syscall.ENOENT {
			return nil, nil, 0, errno
		}
	}
	defer func() {
		if errno == 0 {
			n.root().mu.Lock()
			n.root().files[n.trimPath(p)] = dirty
			n.root().mu.Unlock()
		}
	}()

	fd, err := syscall.Open(p, int(flags)|os.O_CREATE, mode)
	if err != nil {
		return nil, nil, 0, fs.ToErrno(err)
	}

	st := syscall.Stat_t{}
	if err := syscall.Fstat(fd, &st); err != nil {
		syscall.Close(fd)
		return nil, nil, 0, fs.ToErrno(err)
	}

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))
	lf := NewLoopbackFile(fd)

	out.FromStat(&st)
	return ch, lf, 0, 0
}

func (n *loopbackNode) Symlink(ctx context.Context, target, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	err := syscall.Symlink(target, p)
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	st := syscall.Stat_t{}
	if syscall.Lstat(p, &st); err != nil {
		syscall.Unlink(p)
		return nil, fs.ToErrno(err)
	}
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	out.Attr.FromStat(&st)
	return ch, 0
}

func (n *loopbackNode) Link(ctx context.Context, target fs.InodeEmbedder, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	p := filepath.Join(n.path(), name)
	targetNode := toLoopbackNode(target)
	err := syscall.Link(targetNode.path(), p)
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	st := syscall.Stat_t{}
	if syscall.Lstat(p, &st); err != nil {
		syscall.Unlink(p)
		return nil, fs.ToErrno(err)
	}
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	out.Attr.FromStat(&st)
	return ch, 0
}

func (n *loopbackNode) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	p := n.path()

	for l := 256; ; l *= 2 {
		buf := make([]byte, l)
		sz, err := syscall.Readlink(p, buf)
		if err != nil {
			return nil, fs.ToErrno(err)
		}

		if sz < len(buf) {
			return buf[:sz], 0
		}
	}
}

func (n *loopbackNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	p := n.path()
	state := full
	if isWrite(flags) {
		state = dirty
	}
	if err := n.download(p, state); err != nil {
		errno = toErrno(err)
		if !isCreate(flags) || errno != syscall.ENOENT {
			return nil, 0, errno
		}
	}
	if isCreate(flags) {
		defer func() {
			if errno == 0 {
				n.root().mu.Lock()
				n.root().files[n.trimPath(p)] = dirty
				n.root().mu.Unlock()
			}
		}()
	}
	f, err := syscall.Open(p, int(flags), 0)
	if err != nil {
		return nil, 0, fs.ToErrno(err)
	}
	lf := NewLoopbackFile(f)
	return lf, 0, 0
}

func (n *loopbackNode) Opendir(ctx context.Context) syscall.Errno {
	if err := n.download(n.path(), meta); err != nil {
		return toErrno(err)
	}
	fd, err := syscall.Open(n.path(), syscall.O_DIRECTORY, 0755)
	if err != nil {
		return fs.ToErrno(err)
	}
	syscall.Close(fd)
	return fs.OK
}

func (n *loopbackNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	if err := n.download(n.path(), meta); err != nil {
		return nil, toErrno(err)
	}
	return fs.NewLoopbackDirStream(n.path())
}

func (n *loopbackNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if f != nil {
		return f.(fs.FileGetattrer).Getattr(ctx, out)
	}
	p := n.path()
	if err := n.download(p, meta); err != nil {
		return toErrno(err)
	}

	var err error = nil
	st := syscall.Stat_t{}
	err = syscall.Lstat(p, &st)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStat(&st)
	return fs.OK
}

var _ = (fs.NodeSetattrer)((*loopbackNode)(nil))

func (n *loopbackNode) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	p := n.path()
	fsa, ok := f.(fs.FileSetattrer)
	if ok && fsa != nil {
		fsa.Setattr(ctx, in, out)
	} else {
		if m, ok := in.GetMode(); ok {
			if err := syscall.Chmod(p, m); err != nil {
				return fs.ToErrno(err)
			}
		}

		uid, uok := in.GetUID()
		gid, gok := in.GetGID()
		if uok || gok {
			suid := -1
			sgid := -1
			if uok {
				suid = int(uid)
			}
			if gok {
				sgid = int(gid)
			}
			if err := syscall.Chown(p, suid, sgid); err != nil {
				return fs.ToErrno(err)
			}
		}

		mtime, mok := in.GetMTime()
		atime, aok := in.GetATime()

		if mok || aok {

			ap := &atime
			mp := &mtime
			if !aok {
				ap = nil
			}
			if !mok {
				mp = nil
			}
			var ts [2]syscall.Timespec
			ts[0] = fuse.UtimeToTimespec(ap)
			ts[1] = fuse.UtimeToTimespec(mp)

			if err := syscall.UtimesNano(p, ts[:]); err != nil {
				return fs.ToErrno(err)
			}
		}

		if sz, ok := in.GetSize(); ok {
			if err := syscall.Truncate(p, int64(sz)); err != nil {
				return fs.ToErrno(err)
			}
		}
	}

	fga, ok := f.(fs.FileGetattrer)
	if ok && fga != nil {
		fga.Getattr(ctx, out)
	} else {
		st := syscall.Stat_t{}
		err := syscall.Lstat(p, &st)
		if err != nil {
			return fs.ToErrno(err)
		}
		out.FromStat(&st)
	}
	return fs.OK
}

// NewLoopback returns a root node for a loopback file system whose
// root is at the given root. This node implements all NodeXxxxer
// operations available.
func NewLoopbackRoot(root string, c *client.APIClient, opts *Options) (*loopbackRoot, error) {
	var st syscall.Stat_t
	err := syscall.Stat(root, &st)
	if err != nil {
		return nil, err
	}

	n := &loopbackRoot{
		rootPath: root,
		rootDev:  uint64(st.Dev),
		c:        c,
		branches: opts.getBranches(),
		commits:  make(map[string]string),
		files:    make(map[string]fileState),
	}
	return n, nil
}

// download files into the loopback filesystem, if meta is true then only the
// directory structure will be created, no actual data will be downloaded,
// files will be truncated to their actual sizes (but will be all zeros).
func (n *loopbackNode) download(path string, state fileState) (retErr error) {
	path = n.trimPath(path)
	fmt.Println("download: ", path)
	parts := strings.Split(path, "/")
	n.root().mu.Lock()
	cached := n.root().files[path] >= state
	n.root().mu.Unlock()
	if cached {
		fmt.Println("cached")
		return nil
	}
	defer func() {
		if retErr == nil {
			n.root().mu.Lock()
			defer n.root().mu.Unlock()
			n.root().files[path] = state
		}
	}()
	switch {
	case len(parts) == 1 && parts[0] == "":
		ris, err := n.c().ListRepo()
		if err != nil {
			return err
		}
		for _, ri := range ris {
			p := n.repoPath(ri)
			if err := os.MkdirAll(p, 0777); err != nil {
				return err
			}
		}
	default:
		// Make sure the directory for the repo exists
		if err := os.MkdirAll(filepath.Join(n.root().rootPath, parts[0]), 0777); err != nil {
			return err
		}
		commit, err := n.commit(parts[0])
		if err != nil {
			return err
		}
		if commit == "" {
			return nil
		}
		if err := n.c().ListFileF(parts[0], commit, pathpkg.Join(parts[1:]...), 0,
			func(fi *pfs.FileInfo) (retErr error) {
				if fi.FileType == pfs.FileType_DIR {
					return os.MkdirAll(n.filePath(fi), 0777)
				}
				p := n.filePath(fi)
				// Make sure the directory exists
				// I think this may be unnecessary based on the
				if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
					return err
				}
				f, err := os.Create(p)
				if err != nil {
					return err
				}
				defer func() {
					if err := f.Close(); err != nil && retErr == nil {
						retErr = err
					}
				}()
				if state < full {
					return f.Truncate(int64(fi.SizeBytes))
				}
				return n.c().GetFile(fi.File.Commit.Repo.Name, fi.File.Commit.ID, fi.File.Path, 0, 0, f)
			}); err != nil {
			return err
		}
	}
	return nil
}

func (n *loopbackNode) trimPath(path string) string {
	path = strings.TrimPrefix(path, n.root().rootPath)
	return strings.TrimPrefix(path, "/")
}

func (n *loopbackNode) branch(repo string) string {
	if branch, ok := n.root().branches[repo]; ok {
		return branch
	}
	return "master"
}

func (n *loopbackNode) commit(repo string) (string, error) {
	if commit, ok := n.root().commits[repo]; ok {
		return commit, nil
	}
	branch := n.root().branch(repo)
	bi, err := n.root().c.InspectBranch(repo, branch)
	if err != nil && !errutil.IsNotFoundError(err) {
		return "", err
	}
	// You can access branches that don't exist, which allows you to create
	// branches through the fuse mount.
	if errutil.IsNotFoundError(err) || bi.Head == nil {
		n.root().commits[repo] = ""
		return "", nil
	}
	n.root().commits[repo] = bi.Head.ID
	return bi.Head.ID, nil
}

func (n *loopbackNode) repoPath(ri *pfs.RepoInfo) string {
	return filepath.Join(n.root().rootPath, ri.Repo.Name)
}

func (n *loopbackNode) filePath(fi *pfs.FileInfo) string {
	return filepath.Join(n.root().rootPath, fi.File.Commit.Repo.Name, fi.File.Path)
}

func toErrno(err error) syscall.Errno {
	if errutil.IsNotFoundError(err) {
		return syscall.ENOENT
	}
	return fs.ToErrno(err)
}

func isWrite(flags uint32) bool {
	return (int(flags) & (os.O_WRONLY | os.O_RDWR)) != 0
}

func isCreate(flags uint32) bool {
	return int(flags)&os.O_CREATE != 0
}
