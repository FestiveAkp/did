package main

import "github.com/charmbracelet/lipgloss"

var (
	selectedItemStyle  = lipgloss.NewStyle().Background(lipgloss.Color("237"))
	secondaryTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	activeTabStyle     = lipgloss.NewStyle().Bold(true)
	inactiveTabStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	logTimestampStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("71"))
	logSeparatorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	logNoteStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)
