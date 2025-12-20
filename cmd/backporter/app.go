package main

import (
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/backport"
	"codefloe.com/pat-s/backporter/cli/common"
	"codefloe.com/pat-s/backporter/cli/list"
	"codefloe.com/pat-s/backporter/shared/version"
)

func newApp() *cli.Command {
	app := &cli.Command{}
	app.Name = "backporter"
	app.Description = "A tool for backporting git commits and pull requests"
	app.Version = version.String()
	app.Usage = "backport commits and PRs to target branches"
	app.Flags = common.GlobalFlags
	app.Before = common.Before
	app.Suggest = true
	app.Commands = []*cli.Command{
		backport.Command,
		list.Command,
	}

	// Default action when called without subcommand (interactive mode).
	app.Action = backport.Interactive

	return app
}
