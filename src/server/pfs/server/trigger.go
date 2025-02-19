package server

import (
	"time"

	units "github.com/docker/go-units"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/robfig/cron"

	col "github.com/pachyderm/pachyderm/v2/src/internal/collection"
	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/pfsdb"
	"github.com/pachyderm/pachyderm/v2/src/internal/transactionenv/txncontext"
	"github.com/pachyderm/pachyderm/v2/src/pfs"
)

// triggerCommit is called when a commit is finished, it updates branches in
// the repo if they trigger on the change and returns all branches which were
// moved by this call.
func (d *driver) triggerCommit(
	txnCtx *txncontext.TransactionContext,
	commit *pfs.Commit,
) ([]*pfs.Branch, error) {
	repoInfo := &pfs.RepoInfo{}
	if err := d.repos.ReadWrite(txnCtx.SqlTx).Get(pfsdb.RepoKey(commit.Branch.Repo), repoInfo); err != nil {
		return nil, err
	}
	newHead := &pfs.CommitInfo{}
	if err := d.commits.ReadWrite(txnCtx.SqlTx).Get(pfsdb.CommitKey(commit), newHead); err != nil {
		return nil, err
	}
	// find which branches this commit is the head of
	headBranches := make(map[string]bool)
	// The commit _was_ the branch head at some point in time, although it is
	// not guaranteed to still be the branch head. This can happen on a
	// downstream pipeline with triggers - the upstream pipeline may have
	// multiple unfinished commits in its output branch that will be finished
	// one at a time. Without this code, only finishing the _last_ commit would
	// have a chance of triggering.
	headBranches[newHead.Commit.Branch.Name] = true
	for _, b := range repoInfo.Branches {
		bi := &pfs.BranchInfo{}
		if err := d.branches.ReadWrite(txnCtx.SqlTx).Get(pfsdb.BranchKey(b), bi); err != nil {
			return nil, err
		}
		if bi.Head != nil && bi.Head.ID == commit.ID {
			headBranches[b.Name] = true
		}
	}
	triggeredBranches := map[string]bool{}
	var result []*pfs.Branch
	var triggerBranch func(branch *pfs.Branch) error
	triggerBranch = func(branch *pfs.Branch) error {
		if triggeredBranches[branch.Name] {
			return nil
		}
		triggeredBranches[branch.Name] = true
		bi := &pfs.BranchInfo{}
		if err := d.branches.ReadWrite(txnCtx.SqlTx).Get(pfsdb.BranchKey(branch), bi); err != nil {
			return err
		}
		if bi.Trigger != nil {
			if err := triggerBranch(commit.Branch.Repo.NewBranch(bi.Trigger.Branch)); err != nil && !col.IsErrNotFound(err) {
				return err
			}
			if headBranches[bi.Trigger.Branch] {
				var oldHead *pfs.CommitInfo
				if bi.Head != nil {
					oldHead = &pfs.CommitInfo{}
					if err := d.commits.ReadWrite(txnCtx.SqlTx).Get(pfsdb.CommitKey(bi.Head), oldHead); err != nil {
						return err
					}
				}
				triggered, err := d.isTriggered(txnCtx, bi.Trigger, oldHead, newHead)
				if err != nil {
					return err
				}
				if triggered {
					if err := d.branches.ReadWrite(txnCtx.SqlTx).Update(pfsdb.BranchKey(bi.Branch), bi, func() error {
						bi.Head = newHead.Commit
						return nil
					}); err != nil {
						return err
					}
					result = append(result, proto.Clone(commit.Branch).(*pfs.Branch))
					headBranches[bi.Branch.Name] = true
				}
			}
		}
		return nil
	}
	for _, b := range repoInfo.Branches {
		if err := triggerBranch(b); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// isTriggered checks to see if a branch should be updated from oldHead to
// newHead based on a trigger.
func (d *driver) isTriggered(txnCtx *txncontext.TransactionContext, t *pfs.Trigger, oldHead, newHead *pfs.CommitInfo) (bool, error) {
	result := t.All
	merge := func(cond bool) {
		if t.All {
			result = result && cond
		} else {
			result = result || cond
		}
	}
	if t.Size_ != "" {
		size, err := units.FromHumanSize(t.Size_)
		if err != nil {
			// Shouldn't be possible to error here since we validate on ingress
			return false, errors.EnsureStack(err)
		}
		var oldSize uint64
		if oldHead != nil {
			oldSize = oldHead.SizeBytes
		}
		merge(int64(newHead.SizeBytes-oldSize) >= size)
	}
	if t.CronSpec != "" {
		// Shouldn't be possible to error here since we validate on ingress
		schedule, err := cron.ParseStandard(t.CronSpec)
		if err != nil {
			// Shouldn't be possible to error here since we validate on ingress
			return false, errors.EnsureStack(err)
		}
		var oldTime, newTime time.Time
		if oldHead != nil && oldHead.Finished != nil {
			oldTime, err = types.TimestampFromProto(oldHead.Finished)
			if err != nil {
				return false, errors.EnsureStack(err)
			}
		}
		if newHead.Finished != nil {
			newTime, err = types.TimestampFromProto(newHead.Finished)
			if err != nil {
				return false, errors.EnsureStack(err)
			}
		}
		merge(schedule.Next(oldTime).Before(newTime))
	}
	if t.Commits != 0 {
		ci := newHead
		var commits int64
		for commits < t.Commits {
			commits++
			if ci.ParentCommit != nil && (oldHead == nil || oldHead.Commit.ID != ci.ParentCommit.ID) {
				var err error
				ci, err = d.inspectCommit(txnCtx.ClientContext, ci.ParentCommit, pfs.CommitState_STARTED)
				if err != nil {
					return false, err
				}
			} else {
				break
			}
		}
		merge(commits == t.Commits)
	}
	return result, nil
}

// validateTrigger returns an error if a trigger is invalid
func (d *driver) validateTrigger(txnCtx *txncontext.TransactionContext, branch *pfs.Branch, trigger *pfs.Trigger) error {
	if trigger == nil {
		return nil
	}
	if trigger.Branch == "" {
		return errors.Errorf("triggers must specify a branch to trigger on")
	}
	if _, err := cron.ParseStandard(trigger.CronSpec); trigger.CronSpec != "" && err != nil {
		return errors.Wrapf(err, "invalid trigger cron spec")
	}
	if _, err := units.FromHumanSize(trigger.Size_); trigger.Size_ != "" && err != nil {
		return errors.Wrapf(err, "invalid trigger size")
	}
	if trigger.Commits < 0 {
		return errors.Errorf("can't trigger on a negative number of commits")
	}
	bis, err := d.listBranch(txnCtx.ClientContext, branch.Repo, false)
	if err != nil {
		return err
	}
	biMaps := make(map[string]*pfs.BranchInfo)
	for _, bi := range bis {
		biMaps[bi.Branch.Name] = bi
	}
	b := trigger.Branch
	for {
		if b == branch.Name {
			return errors.Errorf("triggers cannot form a loop")
		}
		bi, ok := biMaps[b]
		if ok && bi.Trigger != nil {
			b = bi.Trigger.Branch
		} else {
			break
		}
	}
	return nil
}
