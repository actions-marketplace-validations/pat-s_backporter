// Package list provides the list command for showing backport history.
package list

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/internal"
)

const shaTruncateLength = 12

// Command is the list command.
var Command = &cli.Command{
	Name:   "list",
	Usage:  "list backported commits/PRs from cache",
	Action: listBackports,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "clear",
			Usage: "clear the cache",
		},
	},
}

func listBackports(ctx context.Context, c *cli.Command) error {
	service, err := internal.CreateService(ctx, c)
	if err != nil {
		return err
	}

	if c.Bool("clear") {
		if err := service.ClearCache(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
		fmt.Println("Cache cleared")
		return nil
	}

	entries := service.ListBackports()
	if len(entries) == 0 {
		fmt.Println("No backports found in cache")
		return nil
	}

	fmt.Printf("%-12s %-12s %-20s %-10s %s\n", "ORIGINAL", "BACKPORT", "BRANCH", "PR", "TIMESTAMP")
	fmt.Println("--------------------------------------------------------------------------------------------")

	for _, entry := range entries {
		prStr := "-"
		if entry.PRNumber > 0 {
			prStr = fmt.Sprintf("#%d", entry.PRNumber)
		}

		fmt.Printf("%-12s %-12s %-20s %-10s %s\n",
			safeTruncate(entry.OriginalSHA, shaTruncateLength),
			safeTruncate(entry.BackportSHA, shaTruncateLength),
			entry.TargetBranch,
			prStr,
			entry.Timestamp.Format("2006-01-02 15:04"),
		)
	}

	return nil
}

func safeTruncate(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}
