package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// define initial model
func singleSelectionInitialModel(cfg Configuration) selectionModel {
	treeChoices := []string{}
	for tree := range cfg.Trees {
		treeChoices = append(treeChoices, tree)
	}
	return selectionModel{
		choices:     treeChoices,
		selected:    make(map[int]struct{}),
		multiSelect: false,
		quit:        false,
	}
}

func multiTreeAssignInitialModel(cfg Configuration, proposedState State) selectionModel {

	// all trees
	treeChoices := []string{}
	for tree := range cfg.Trees {
		treeChoices = append(treeChoices, tree)
	}

	// trees where branch is already linked
	linkedTrees := make(map[int]struct{})
	for tree, states := range cfg.Trees {
		for _, state := range states.States {
			if state.Branch == proposedState.Branch && state.Repo == proposedState.Repo {
				linkedTrees[getIndex(treeChoices, tree)] = struct{}{}
			}
		}
	}

	repoChoices := []string{}
	for repo := range cfg.Repositories {
		repoChoices = append(repoChoices, repo)
	}
	return selectionModel{
		choices:     treeChoices,
		selected:    linkedTrees,
		multiSelect: true,
		quit:        false,
	}
}

// run nothing on init
func (m selectionModel) Init() tea.Cmd {
	return nil
}

// update model
func (m selectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// clear selection
			m.quit = true
			m.selected = make(map[int]struct{})
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "s":
			if m.multiSelect {
				return m, tea.Quit
			}
			fmt.Println("\nYour selected choices:")
		case " ":
			_, ok := m.selected[m.cursor]
			if m.multiSelect {
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectionModel) View() string {
	// The header
	s := "Select a tree to load:\n\n"
	if m.multiSelect {
		s = "Select/unselect the trees to link with this branch:\n\n"
	}

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " "
		if m.cursor == i {
			cursor = "→"
		}

		// Is this choice selected?
		checked := "□" // not selected
		if _, ok := m.selected[i]; ok {
			checked = "■" // selected!
		}

		// Render the row
		if m.multiSelect {
			s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}
	if m.multiSelect {
		s += fmt.Sprint("\n\033[31m", "space = Select   enter = Save   q = Quit", "\n\033[0m")
	} else {

		s += fmt.Sprint("\n\033[31m", "enter = Select    q = Quit", "\n\033[0m")
	}
	return s
}
