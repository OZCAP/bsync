package main

import (
	"os"
	"fmt"
	"github.com/urfave/cli/v2"
	"bufio"
	"strings"
)

func newCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		// terminal input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(enterTextStyle.Render("Enter name for new tree: \u001b[31;1m"))
		newTreeName, _ := reader.ReadString('\n')
		fmt.Print("\u001b[0m")
		fmt.Println(newTreeStyle.Render(fmt.Sprint("🌳 Created new tree: ", newTreeName)))

		newTreeName = strings.TrimSpace(newTreeName)
		newTreeOwner := getGitUser()

		// create new tree
		newTree := Tree{
			Name:   newTreeName,
			States: []State{},
			Owner:  newTreeOwner,
		}

		// if config.Trees doesn't exist, create it
		if config.Trees == nil {
			config.Trees = make(map[string]Tree)
		}

		config.Trees[newTreeName] = newTree
		saveConfigToml(config)
		return nil
	}
}
