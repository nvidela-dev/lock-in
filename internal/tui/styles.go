package tui

import "charm.land/lipgloss/v2"

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(true)
	doneStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("108"))
	readyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("222"))
	progressStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	statusStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("108"))
	footerStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	projectActive     = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("24")).Bold(true)
	projectInactive   = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	manualBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)
