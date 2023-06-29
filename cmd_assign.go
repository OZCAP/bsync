package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"os"
)

func assignCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {

		// assign branch to trees
		branch := getBranchName()
		remoteRepository := getRemoteRepository()
		repositoryName := parseRepositoryName(remoteRepository)

		proposedState := State{
			Branch: branch,
			Repo:   repositoryName,
		}

		p := tea.NewProgram(multiTreeAssignInitialModel(config, proposedState))
		finalModel, err := p.StartReturningModel()
		if err != nil {
			fmt.Println("Oh no, it broke!")
			os.Exit(1)
		}
		if m, ok := finalModel.(selectionModel); ok {
			if m.quit {
				return nil
			}
			newTrees := makeNewTreesFromSelection(m, branch, repositoryName, config)
			config.Trees = newTrees
			saveConfigToml(config)
		}
		return nil
	}
}
