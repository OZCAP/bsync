package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
	"bufio"
	"os"
)

func pullRequestCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		firstArg := cCtx.Args().Get(0)
		pullBranches := cCtx.Args().Tail()

		if firstArg != "" {
			pullBranches = append([]string{firstArg}, pullBranches...)
		}

		if len(pullBranches) == 0 {
			fmt.Println("Please specify a branch to open a pull request for")
			return nil
		}

		branchName := getBranchName()

		// prompt user to input Monday ticket id (if applicable)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Monday ticket ID (default: #): ")
		ticketID, _ := reader.ReadString('\n')
		ticketID = strings.TrimSpace(ticketID) // remove the newline
		if !strings.HasPrefix(ticketID, "#") {
			ticketID = "#" + ticketID
		}

		// open pull requests for each destination branch
		for _, destinationBranch := range pullBranches {
			pullTitle := formatPRTitle(branchName, destinationBranch)
			openPullRequest(pullTitle, destinationBranch, ticketID)
		}

		return nil
	}
}
