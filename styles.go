package main

import (
	"github.com/charmbracelet/lipgloss"
)

var newTreeStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#EF4160")).
	Background(lipgloss.Color("#fcfcfc")).
	PaddingTop(1).
	PaddingLeft(3).
	PaddingRight(3).
	MarginTop(1).
	MarginBottom(1)


var enterTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7dc088")).
	PaddingTop(1)