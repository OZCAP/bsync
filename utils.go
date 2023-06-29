package main

import (
	"crypto/sha256"
	"fmt"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
	"time"
	"sort"
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
		fmt.Printf("\npulling branch %s 🪵\n", state.Branch)
		cmd := exec.Command("git", "checkout", state.Branch)
		cmd.Dir = cfg.Repositories[state.Repo].Local
		fmt.Printf("\033[F\033[2K") // move cursor up and clear line
		_, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error checking out branch %s (%s) ❌", state.Branch, state.Repo)
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Checked out branch %s (%s) ✅\n", state.Branch, state.Repo)
			fmt.Printf("Pulling branch %s (%s) 🪵\n", state.Branch, state.Repo)
			pullCmd := exec.Command("git", "pull", "--ff-only", "origin", state.Branch)
			pullCmd.Dir = cfg.Repositories[state.Repo].Local
			_, err :=pullCmd.Output()
			if err != nil {
				fmt.Printf(err.Error())
				fmt.Printf("\033[F\033[2KError pulling branch %s (%s) ❌\n", state.Branch, state.Repo)
			} else {
				fmt.Printf("\033[F\033[2KPulled branch %s (%s) ✅\n", state.Branch, state.Repo)
			}
		}
		fmt.Println("")
	}
	// fmt.Println("\033[2KLoaded tree", project, "🌳")
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

func openPullRequest(title string, destinationBranch string, body string) {
    cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", body, "--base", destinationBranch)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    fmt.Println("\u001b[31mCreating pull request into branch " + destinationBranch + "...\u001b[0m" + "\n\u001b[34m - " + title + "\u001b")
    cmd.Run()
    fmt.Println()
    return
}

type BranchTime struct {
	Branch string
	Time   time.Time
}

func listBranches() {
	// grep a list of all branches in local git repo
	// cmd := exec.Command("git", "branch", "--list")
	// out, err := cmd.Output()
	// if err != nil {
	// 	panic(err)
	// }

	// output := string(out)

	// // remove asterisk from output
	// output = strings.ReplaceAll(output, "*", "")

	// // split branches into array
	// branches := strings.Split(string(out), "\n")

	// for i, branch := range branches {
	// 	branches[i] = strings.TrimSpace(branch)
	// 	fmt.Println(i, branch)
	// }
allBranches := []string{}
branchesWithTime := []BranchTime{}
subDirs := []string{}
	// list files in .git/refs/heads
	items, err := os.ReadDir(".git/refs/heads")
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		// open dir if it is and add to list
		if item.IsDir() {
			subDirs = append(subDirs, item.Name())
		}
	}

	for _, subDir := range subDirs {
		items, err := os.ReadDir(".git/refs/heads/" + subDir)
		if err != nil {
			panic(err)
		}

		for _, item := range items {
			allBranches = append(allBranches, subDir + "/" + item.Name())
		}
	}



	for _, branch := range allBranches {
		// sort branches by last commit date
		cmd := exec.Command("git", "log", "-1", "--format=%cd", branch)
		out, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		// convert date to unix timestamp
		date, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", string(out))
		if err != nil {
			panic(err)
		}

		fmt.Println(branch, date)

		// add branch and time to array
		branchesWithTime = append(branchesWithTime, BranchTime{
			Branch: branch,
			Time:   date,
		})
	}

	fmt.Println(branchesWithTime)

	// sort branches by time
	sort.Slice(branchesWithTime, func(i, j int) bool {
		return branchesWithTime[i].Time.After(branchesWithTime[j].Time)
	})

	// print branches
	for _, branch := range branchesWithTime {
		fmt.Println(branch.Branch)
	}






	// for _, file := range files {
	// 	fmt.Println(file.Name())
	// }

}