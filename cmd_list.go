package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func listCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		for _, tree := range config.Trees {
			fmt.Println(tree.Name)
			for _, state := range tree.States {
				fmt.Println("⊢", state.Repo, "("+state.Branch+")")
			}
		}
		return nil
	}
}
