package views

import "github.com/charmbracelet/lipgloss"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

var infoMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#489eff"))
var warningMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffed68"))
var errorMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f94848"))
