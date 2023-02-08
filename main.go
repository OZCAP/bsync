package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
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
				Name:  "load",
				Usage: "load tree and pull branches",
				Action: func(cCtx *cli.Context) error {
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
				},
			}, {
				Name:    "add",
				Aliases: []string{"add-repo"},
				Usage:   "add repository to local config",
				Action: func(cCtx *cli.Context) error {
					// add repository to config
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
				},
			}, {
				Name:    "new",
				Aliases: []string{"new-tree"},
				Usage:   "create new tree project",
				Action: func(cCtx *cli.Context) error {
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
				},
			}, {
				Name:    "assign",
				Aliases: []string{"assign-branch"},
				Usage:   "assign current branch to tree",
				Action: func(cCtx *cli.Context) error {

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
				},
			}, {
				Name:    "ls",
				Aliases: []string{"list"},
				Usage:   "list trees",
				Action: func(cCtx *cli.Context) error {
					for _, tree := range config.Trees {
						fmt.Println(tree.Name)
						for _, state := range tree.States {
							fmt.Println("⊢", state.Repo, "("+state.Branch+")")
						}
					}
					return nil
				},
			}, {
				Name:    "rm",
				Aliases: []string{"remove", "delete"},
				Usage:   "remove tree",
				Action: func(cCtx *cli.Context) error {
					treeNameToRemove := cCtx.Args().Get(0)
					if treeNameToRemove == "" {
						fmt.Println("Please specify a tree to remove")
						return nil
					}
					delete(config.Trees, treeNameToRemove)
					saveConfigToml(config)
					return nil
				},
				}, {
					Name:    "pr",
					Aliases: []string{"pull-request"},
					Usage:   "open pull requests for current branch",
					Action: func(cCtx *cli.Context) error {
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

						// open pull requests for each destination branch
						for _, destinationBranch := range pullBranches {
							pullTitle := formatPRTitle(branchName, destinationBranch)
							openPullRequest(pullTitle, destinationBranch)
						}

						return nil
					},
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
