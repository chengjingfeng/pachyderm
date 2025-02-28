package pfsload

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"path"

	"github.com/pachyderm/pachyderm/v2/src/internal/randutil"
)

const pathSize = 32

type MemFile struct {
	path    string
	content []byte
}

func NewMemFile(path string, data []byte) *MemFile {
	return &MemFile{
		path:    path,
		content: data,
	}
}

func (mf *MemFile) Path() string {
	return mf.path
}

func (mf *MemFile) Reader() io.Reader {
	return bytes.NewReader(mf.content)
}

type FileSourceSpec struct {
	Name                 string                `yaml:"name,omitempty"`
	RandomFileSourceSpec *RandomFileSourceSpec `yaml:"random,omitempty"`
}

type FileSource interface {
	Next() (*MemFile, error)
}

func NewFileSource(spec *FileSourceSpec, random *rand.Rand) FileSource {
	return newRandomFileSource(spec.RandomFileSourceSpec, random)
}

type RandomFileSourceSpec struct {
	IncrementPath       bool                 `yaml:"incrementPath,omitempty"`
	RandomDirectorySpec *RandomDirectorySpec `yaml:"directory,omitempty"`
	SizeSpecs           []*SizeSpec          `yaml:"size,omitempty"`
}

type randomFileSource struct {
	spec      *RandomFileSourceSpec
	random    *rand.Rand
	dirSource *randomDirectorySource
	next      int64
}

func newRandomFileSource(spec *RandomFileSourceSpec, random *rand.Rand) FileSource {
	var dirSource *randomDirectorySource
	if spec.RandomDirectorySpec != nil {
		dirSource = &randomDirectorySource{
			spec:   spec.RandomDirectorySpec,
			random: random,
		}
	}
	return &randomFileSource{
		spec:      spec,
		random:    random,
		dirSource: dirSource,
	}
}

func (rfs *randomFileSource) Next() (*MemFile, error) {
	sizeSpec, err := FuzzSize(rfs.spec.SizeSpecs, rfs.random)
	if err != nil {
		return nil, err
	}
	min, max := sizeSpec.Min, sizeSpec.Max
	size := min
	if max > min {
		size += rfs.random.Intn(max - min)
	}
	return NewMemFile(rfs.nextPath(), randutil.Bytes(rfs.random, size)), nil
}

func (rfs *randomFileSource) nextPath() string {
	var dir string
	if rfs.dirSource != nil {
		dir = rfs.dirSource.nextPath()
	}
	if rfs.spec.IncrementPath {
		next := rfs.next
		rfs.next += 1
		return path.Join(dir, fmt.Sprintf("%016d", next))
	}
	return path.Join(dir, string(randutil.Bytes(rfs.random, pathSize)))
}

type RandomDirectorySpec struct {
	Depth int   `yaml:"depth,omitempty"`
	Run   int64 `yaml:"run,omitempty"`
}

type randomDirectorySource struct {
	spec   *RandomDirectorySpec
	random *rand.Rand
	next   string
	run    int64
}

func (rds *randomDirectorySource) nextPath() string {
	if rds.next == "" {
		depth := rds.random.Intn(rds.spec.Depth)
		for i := 0; i < depth; i++ {
			rds.next = path.Join(rds.next, string(randutil.Bytes(rds.random, pathSize)))
		}
	}
	dir := rds.next
	rds.run++
	if rds.run == rds.spec.Run {
		rds.next = ""
		rds.run = 0
	}
	return dir

}

type FilesSpec struct {
	Count     int         `yaml:"count,omitempty"`
	FileSpecs []*FileSpec `yaml:"file,omitempty"`
}

func Files(env *Env, spec *FilesSpec) ([]*MemFile, error) {
	var files []*MemFile
	for i := 0; i < spec.Count; i++ {
		file, err := FuzzFile(env, spec.FileSpecs)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// TODO: Add different types of files.
type FileSpec struct {
	Source string `yaml:"source,omitempty"`
	Prob   int    `yaml:"prob,omitempty"`
}

func File(env *Env, spec *FileSpec) (*MemFile, error) {
	return env.FileSource(spec.Source).Next()
}

type SizeSpec struct {
	Min  int `yaml:"min,omitempty"`
	Max  int `yaml:"max,omitempty"`
	Prob int `yaml:"prob,omitempty"`
}
