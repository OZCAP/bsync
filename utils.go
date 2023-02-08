package main

import (
	"crypto/sha256"
	"fmt"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

// load config from toml
func loadConfigToml() Configuration {

	// define config location
	homeDir, e := os.UserHomeDir()
	if e != nil {
		panic(e)
	}
	configPath := homeDir + "/Library/Application Support/bsync"
	configFileName := "config.toml"
	fullConfigPath := configPath + "/" + configFileName

	// check if config file exists and create if not
	if _, err := os.Stat(configPath + "/" + configFileName); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0700)
		os.WriteFile(fullConfigPath, []byte{}, 0644)
	}

	// load config file data
	doc, e := os.ReadFile(fullConfigPath)
	if e != nil {
		panic(e)
	}

	// deserialize config file
	var cfg Configuration
	err := toml.Unmarshal(doc, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

// get working directory of repository
func getLocalRepository() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	currentRepository := strings.TrimSpace(string(out))
	return currentRepository
}

// get remote repository url
func getRemoteRepository() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	repositoryUrl := strings.TrimSpace(string(out))
	return repositoryUrl
}

func parseRepositoryName(repositoryUrl string) string {

	if strings.HasSuffix(repositoryUrl, ".git") {
		repositoryUrl = repositoryUrl[:len(repositoryUrl)-4]
	}

	repositoryName := strings.Split(repositoryUrl, ".com/")[1]

	return repositoryName
}

func getBranchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	currentBranch := strings.TrimSpace(string(out))
	return currentBranch
}

func loadProject(project string, cCtx *cli.Context, cfg Configuration) {
	projectDetails := cfg.Trees[project]
	// fmt.Println("trees:", cfg.Trees, "project:", project)
	projectStates := projectDetails.States
	for _, state := range projectStates {
		fmt.Printf("\npulling branch %s 🪵", state.Branch)

		cmd := exec.Command("git", "checkout", state.Branch)
		cmd.Dir = cfg.Repositories[state.Repo].Local

		fmt.Printf("\033[F\033[2K") // move cursor up and clear line
		_, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error pulling branch %s ❌", state.Branch)
			panic(err)
		} else {
			fmt.Printf("Pulled branch %s (%s) ✅", state.Branch, state.Repo)
		}
		fmt.Println("")
	}
	fmt.Println("\033[2KLoaded tree", project, "🌳")
}

func saveConfigToml(cfg Configuration) {
	// define config location
	homeDir, e := os.UserHomeDir()
	if e != nil {
		panic(e)
	}
	configPath := homeDir + "/Library/Application Support/bsync"
	configFileName := "config.toml"
	fullConfigPath := configPath + "/" + configFileName

	b, err := toml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	os.WriteFile(fullConfigPath, b, 0644)
}

func getIndex(choices []string, choice string) int {
	for i, c := range choices {
		if c == choice {
			return i
		}
	}
	return -1
}

func remove(s []State, i int) []State {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func makeNewTreesFromSelection(m selectionModel, branch string, repositoryName string, config Configuration) map[string]Tree {
	gitUser := getGitUser()

	selectedTreeNames := []string{}
	for _, choice := range m.choices {
		for i, c := range m.choices {
			if c == choice {
				if _, ok := m.selected[i]; ok {
					selectedTreeNames = append(selectedTreeNames, choice)
				}
			}
		}
	}

	unSelectedTreeNames := []string{}
	for _, choice := range m.choices {
		found := false
		for _, selectedTree := range selectedTreeNames {
			if selectedTree == choice {
				found = true
			}
		}
		if !found {
			unSelectedTreeNames = append(unSelectedTreeNames, choice)
		}
	}

	newTrees := map[string]Tree{}
	for _, treeName := range selectedTreeNames {
		tree := config.Trees[treeName]
		treeStates := tree.States

		// check if branch already exists
		branchExists := false
		for _, state := range tree.States {
			if state.Branch == branch && state.Repo == repositoryName {
				branchExists = true
			}
		}

		// remove sibling branch
		for i, state := range treeStates {
			if state.Branch != branch && state.Repo == repositoryName {
				treeStates = remove(treeStates, i)
			}
		}

		// if branch doesn't exist, add it
		if !branchExists {
			treeStates = append(treeStates, State{
				Branch: branch,
				Repo:   repositoryName,
			})
		}

		// update tree
		newTree := Tree{
			Name:   tree.Name,
			Owner:  gitUser,
			States: treeStates,
		}
		newTrees[tree.Name] = newTree
	}

	for _, treeName := range unSelectedTreeNames {
		tree := config.Trees[treeName]
		treeStates := tree.States

		// check if branch already exists
		branchExists := false
		for _, state := range treeStates {
			if state.Branch == branch && state.Repo == repositoryName {
				branchExists = true
			}
		}

		// if branch exists, remove it
		if branchExists {
			for i, state := range treeStates {
				if state.Branch == branch && state.Repo == repositoryName {
					treeStates = remove(treeStates, i)
				}
			}
		}

		// update tree
		newTree := Tree{
			Name:   tree.Name,
			Owner:  gitUser,
			States: treeStates,
		}
		newTrees[tree.Name] = newTree
	}

	return newTrees
}

func getGitUser() string {
	cmd := exec.Command("git", "config", "--get", "user.email")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	user := strings.ToUpper(strings.TrimSpace(string(out)))
	return user
}

type PatternFormat struct {
	Pattern     string
	Replacement string
	Placeholder string
}

func formatPRTitle(currentBranch string, destinationBranch string) string {
	branchNameTitle := currentBranch

	rawPatterns := []PatternFormat{
		{
			Pattern:     "--",
			Replacement: " -- ",
			Placeholder: "",
		},
		{
			Pattern:     "-",
			Replacement: " ",
			Placeholder: "",
		},
	}

	patterns := []PatternFormat{}

	// replace all patterns with placeholders
	for _, pattern := range rawPatterns {
		pattern.Placeholder = fmt.Sprintf("{%x}", sha256.Sum256([]byte(pattern.Pattern)))
		patterns = append(patterns, pattern)
		branchNameTitle = strings.ReplaceAll(branchNameTitle, pattern.Pattern, pattern.Placeholder)
	}

	// replace all placeholders with replacements
	for _, pattern := range patterns {
		branchNameTitle = strings.ReplaceAll(branchNameTitle, pattern.Placeholder, pattern.Replacement)
	}

	prefix := "[" + destinationBranch + "]"
	return prefix + " " + branchNameTitle
}

func openPullRequest(title string, destinationBranch string) {
	cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", "", "--base", destinationBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("\u001b[31mCreating pull request into branch " + destinationBranch + "...\u001b[0m" + "\n\u001b[34m - " + title + "\u001b")
	cmd.Run()
	fmt.Println()
	return
}
