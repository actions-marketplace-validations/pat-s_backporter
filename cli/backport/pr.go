package backport

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/internal"
	"codefloe.com/pat-s/backporter/pkg/backport"
	"codefloe.com/pat-s/backporter/shared/logger"
)

func backportPR(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() < 2 { //nolint:mnd
		return fmt.Errorf("usage: backport pr <pr-number> <target-branch>")
	}

	prNumberStr := c.Args().Get(0)
	targetBranch := c.Args().Get(1)
	dryRun := c.Bool("dry-run")

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		return fmt.Errorf("invalid PR number: %s", prNumberStr)
	}

	service, err := internal.CreateService(ctx, c)
	if err != nil {
		return err
	}

	opts := backport.BackportOptions{
		TargetBranch: targetBranch,
		DryRun:       dryRun,
	}

	result, err := service.BackportPR(ctx, prNumber, opts)
	if err != nil {
		return err
	}

	return handleBackportResult(result)
}

func handleBackportResult(result *backport.BackportResult) error {
	if result.HasConflict {
		log.Debug().Msg("cherry-pick resulted in conflicts")

		if logger.IsCI() {
			return fmt.Errorf("cherry-pick conflicts detected in CI mode")
		}

		fmt.Println()
		fmt.Println("✗ Cherry-pick resulted in conflicts")
		fmt.Println()
		fmt.Println("To resolve:")
		fmt.Println("  1. Fix the conflicts in the affected files")
		fmt.Println("  2. Run: git cherry-pick --continue")
		fmt.Println()
		fmt.Println("To abort:")
		fmt.Println("  Run: git cherry-pick --abort")
		fmt.Println()
		fmt.Println("Conflict details:")
		fmt.Println(result.Message)

		return fmt.Errorf("cherry-pick conflicts need resolution")
	}

	if result.Success {
		log.Debug().
			Str("original", result.OriginalSHA).
			Str("backport", result.BackportSHA).
			Str("branch", result.TargetBranch).
			Msg("backport completed successfully")

		// Pretty output for successful backport.
		shortOriginal := result.OriginalSHA
		if len(shortOriginal) > 8 { //nolint:mnd
			shortOriginal = shortOriginal[:8]
		}
		shortBackport := result.BackportSHA
		if len(shortBackport) > 8 { //nolint:mnd
			shortBackport = shortBackport[:8]
		}

		fmt.Println()
		if result.PRNumber > 0 {
			fmt.Printf("✓ Successfully backported PR #%d to %s\n", result.PRNumber, result.TargetBranch)
		} else {
			fmt.Printf("✓ Successfully backported commit %s to %s\n", shortOriginal, result.TargetBranch)
		}
		fmt.Printf("  New commit: %s\n", shortBackport)
		fmt.Println()
	}

	return nil
}
