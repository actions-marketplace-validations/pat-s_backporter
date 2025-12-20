// Package backport provides CLI commands for backporting.
package backport

import (
	"github.com/urfave/cli/v3"
)

// Command is the root backport command.
var Command = &cli.Command{
	Name:  "backport",
	Usage: "backport commits or PRs to target branches",
	Commands: []*cli.Command{
		prCmd,
		commitCmd,
	},
}

var prCmd = &cli.Command{
	Name:      "pr",
	Usage:     "backport a pull request",
	ArgsUsage: "<pr-number> <target-branch>",
	Action:    backportPR,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show what would be done without making changes",
		},
	},
}

var commitCmd = &cli.Command{
	Name:      "commit",
	Usage:     "backport a commit",
	ArgsUsage: "<commit-sha> <target-branch>",
	Action:    backportCommit,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show what would be done without making changes",
		},
	},
}
