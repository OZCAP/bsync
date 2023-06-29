package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	// load local config
	config := loadConfigToml()

	// run CLI app
	app := &cli.App{
		Name:  "\u001b[31;1mbsync\u001b[0m",
		Usage: "Sync branches across multiple git repositories",
		Commands: []*cli.Command{
			{
				Name:   "load",
				Usage:  "load tree and pull branches",
				Action: loadCmdAction(config),
			}, {
				Name:    "add",
				Aliases: []string{"add-repo"},
				Usage:   "add repository to local config",
				Action:  addCmdAction(config),
			}, {
				Name:    "new",
				Aliases: []string{"new-tree"},
				Usage:   "create new tree project",
				Action:  newCmdAction(config),
			}, {
				Name:    "assign",
				Aliases: []string{"assign-branch"},
				Usage:   "assign current branch to tree",
				Action:  assignCmdAction(config),
			}, {
				Name:    "ls",
				Aliases: []string{"list"},
				Usage:   "list trees",
				Action:  listCmdAction(config),
			}, {
				Name:    "rm",
				Aliases: []string{"remove", "delete"},
				Usage:   "remove tree",
				Action:  removeCmdAction(config),
			}, {
				Name:    "pr",
				Aliases: []string{"pull-request"},
				Usage:   "open pull requests for current branch",
				Action:  pullRequestCmdAction(config),
			},
			{
				Name:    "branch",
				Aliases: []string{"switch-branch"},
				Usage:   "switch to branch in tree",
				Action: branchCmdAction(config),
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "",
				Usage: "name of tree",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
