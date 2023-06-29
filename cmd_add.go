package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func addCmdAction(config Configuration) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		localRepo := getLocalRepository()
		remoteRepo := getRemoteRepository()
		currentRepositoryName := parseRepositoryName(remoteRepo)

		// try and find existing repo
		repoExists := false
		for _, repo := range config.Repositories {
			if repo.Remote == remoteRepo {
				repoExists = true
			}
		}
		if repoExists {
			fmt.Println("Repository\u001b[31;1m", currentRepositoryName, "\u001b[0malready exists")
		} else {
			newRepo := Repository{
				Remote: remoteRepo,
				Local:  localRepo,
			}

			// if config.Repositories doesn't exist, create it
			if config.Repositories == nil {
				config.Repositories = make(map[string]Repository)
			}

			config.Repositories[currentRepositoryName] = newRepo
			saveConfigToml(config)
			fmt.Println("Added repository\u001b[31;1m", currentRepositoryName, "\u001b[0m")
		}

		return nil
	}
}
