package main

import "github.com/charmbracelet/lipgloss"

var (
	selectedItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true)
	secondaryTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	activeTabStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true)
	inactiveTabStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	activeFilterStyle   = lipgloss.NewStyle().Bold(true)
	inactiveFilterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	logTimestampStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("71"))
	logSeparatorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	logNoteStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)
