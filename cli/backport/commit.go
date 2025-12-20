package backport

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/internal"
	"codefloe.com/pat-s/backporter/pkg/backport"
)

func backportCommit(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() < 2 { //nolint:mnd
		return fmt.Errorf("usage: backport commit <commit-sha> <target-branch>")
	}

	sha := c.Args().Get(0)
	targetBranch := c.Args().Get(1)
	dryRun := c.Bool("dry-run")

	service, err := internal.CreateService(ctx, c)
	if err != nil {
		return err
	}

	opts := backport.BackportOptions{
		TargetBranch: targetBranch,
		DryRun:       dryRun,
	}

	result, err := service.BackportCommit(ctx, sha, opts)
	if err != nil {
		return err
	}

	return handleBackportResult(result)
}
