package service

import (
	"context"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/pachyderm/pachyderm/v2/src/client"
	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/ppsutil"
	"github.com/pachyderm/pachyderm/v2/src/internal/uuid"
	"github.com/pachyderm/pachyderm/v2/src/pfs"
	"github.com/pachyderm/pachyderm/v2/src/pps"
	"github.com/pachyderm/pachyderm/v2/src/server/worker/common"
	"github.com/pachyderm/pachyderm/v2/src/server/worker/datum"
	"github.com/pachyderm/pachyderm/v2/src/server/worker/driver"
	"github.com/pachyderm/pachyderm/v2/src/server/worker/logs"
)

// Run will run a service pipeline until the driver is canceled.
// TODO: The context handling is wonky here, the pachClient context is above the service context in the hierarchy.
// This is necessary to ensure we can finish the job when the service gets canceled. Services will probably be reworked to not
// be triggered by output commits, so this is probably fine for now.
func Run(driver driver.Driver, logger logs.TaggedLogger) error {
	pachClient := driver.PachClient()
	pipelineInfo := driver.PipelineInfo()
	return forEachPipelineJob(pachClient, pipelineInfo, logger, func(ctx context.Context, pipelineJobInfo *pps.PipelineJobInfo) (retErr error) {
		driver := driver.WithContext(ctx)
		if err := driver.UpdatePipelineJobState(pipelineJobInfo.PipelineJob.ID, pps.PipelineJobState_JOB_RUNNING, ""); err != nil {
			return err
		}
		pipelineJobInput := ppsutil.PipelineJobInput(pipelineInfo, pipelineJobInfo.OutputCommit)
		di, err := datum.NewIterator(pachClient, pipelineJobInput)
		if err != nil {
			return err
		}
		var meta *datum.Meta
		if err := di.Iterate(func(m *datum.Meta) error {
			if meta != nil {
				return errors.New("services must have a single datum")
			}
			meta = m
			return nil
		}); err != nil {
			return err
		}
		if meta == nil {
			return errors.New("services must have a single datum")
		}
		defer func() {
			if common.IsDone(ctx) {
				retErr = finishPipelineJob(pachClient, pipelineJobInfo)
			}
		}()
		storageRoot := filepath.Join(driver.InputDir(), client.PPSScratchSpace, uuid.NewWithoutDashes())
		return datum.WithSet(pachClient, storageRoot, func(s *datum.Set) error {
			inputs := meta.Inputs
			logger = logger.WithData(inputs)
			env := driver.UserCodeEnv(logger.PipelineJobID(), pipelineJobInfo.OutputCommit, inputs)
			return s.WithDatum(ctx, meta, func(d *datum.Datum) error {
				return driver.WithActiveData(inputs, d.PFSStorageRoot(), func() error {
					return d.Run(ctx, func(runCtx context.Context) error {
						return driver.RunUserCode(runCtx, logger, env)
					})
				})
			})

		})
	})
}

// Repeatedly runs the given callback with the latest pipelineJob for the pipeline.
// The given context will be canceled if a newer commit is ready, then this will
// wait for the previous callback to return before calling the callback again
// with the latest commit.
func forEachPipelineJob(pachClient *client.APIClient, pipelineInfo *pps.PipelineInfo, logger logs.TaggedLogger, cb func(context.Context, *pps.PipelineJobInfo) error) error {
	// These are used to cancel the existing service and wait for it to finish
	var cancel func()
	var eg *errgroup.Group
	return pachClient.SubscribePipelineJob(pipelineInfo.Pipeline.Name, false, func(pji *pps.PipelineJobInfo) error {
			if cancel != nil {
				logger.Logf("canceling previous service, new commit ready")
				cancel()
				if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
					return err
				}
			}
			logger.Logf("starting new service, commit: %s", pji.OutputCommit.ID)
			var ctx context.Context
			ctx, cancel = context.WithCancel(pachClient.Ctx())
			eg, ctx = errgroup.WithContext(ctx)
			eg.Go(func() error { return cb(ctx, pji) })
			return nil
		},
	)
}

func finishPipelineJob(pachClient *client.APIClient, pipelineJobInfo *pps.PipelineJobInfo) error {
	_, err := pachClient.RunBatchInTransaction(func(builder *client.TransactionBuilder) error {
		if _, err := builder.PfsAPIClient.FinishCommit(pachClient.Ctx(), &pfs.FinishCommitRequest{
			Commit: pipelineJobInfo.OutputCommit,
		}); err != nil {
			return err
		}
		_, err := builder.PpsAPIClient.UpdatePipelineJobState(pachClient.Ctx(), &pps.UpdatePipelineJobStateRequest{
			PipelineJob: pipelineJobInfo.PipelineJob,
			State:       pps.PipelineJobState_JOB_SUCCESS,
		})
		return err
	})
	return err
}
