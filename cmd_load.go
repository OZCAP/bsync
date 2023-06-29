package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"os"
)

func loadCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		projectNameToLoad := cCtx.Args().Get(0)
		if projectNameToLoad == "" {

			// select project by ui
			p := tea.NewProgram(singleSelectionInitialModel(config))
			finalModel, err := p.StartReturningModel()
			if err != nil {
				fmt.Println("Oh no, it broke!")
				os.Exit(1)
			}
			if m, ok := finalModel.(selectionModel); ok {
				if m.quit {
					return nil
				}
				fmt.Println("Loading tree:", m.choices[m.cursor])
				projectNameToLoad = m.choices[m.cursor]
			}
		}
		if projectNameToLoad != "" {
			if _, ok := config.Trees[projectNameToLoad]; ok {
				loadProject(projectNameToLoad, cCtx, config)
			} else {
				fmt.Println("Project", projectNameToLoad, "does not exist")
			}

		} else {
			fmt.Println("No project name provided")
		}
		return nil
	}
}
