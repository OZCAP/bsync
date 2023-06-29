package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func removeCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		treeNameToRemove := cCtx.Args().Get(0)
		if treeNameToRemove == "" {
			fmt.Println("Please specify a tree to remove")
			return nil
		}
		delete(config.Trees, treeNameToRemove)
		saveConfigToml(config)
		return nil
	}
}
