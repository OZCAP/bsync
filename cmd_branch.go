package main

import (
	// "fmt"
	"github.com/urfave/cli/v2"
	// "os"
)

func branchCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		listBranches()

		return nil

	}
}
