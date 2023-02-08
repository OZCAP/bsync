package main

type Configuration struct {
	Trees        map[string]Tree       `toml:"trees"`
	Repositories map[string]Repository `toml:"repositories"`
	ActiveTree   string                `toml:"active_tree"`
}

type Tree struct {
	Name   string  `toml:"name"`
	Owner  string  `toml:"owner"`
	States []State `toml:"states"`
}

type Repository struct {
	Remote string `toml:"remote"`
	Local  string `toml:"local"`
}

type selectionModel struct {
	choices     []string
	cursor      int
	selected    map[int]struct{}
	multiSelect bool
	quit        bool
}

type State struct {
	Repo   string `toml:"repo"`
	Branch string `toml:"branch"`
}
