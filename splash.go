package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type splashModel struct {
	ready bool
}

func initialSplashModel() splashModel {
	return splashModel{}
}

func (m splashModel) Init() tea.Cmd {
	return nil
}

func (m splashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			m.ready = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m splashModel) View() string {
	if m.ready {
		return ""
	}

	// Define styles
	// titleStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("#FF6B6B")).
	// 	Bold(true).
	// 	Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Align(lipgloss.Center)

	// Create the logo
	logo := `
███████╗ █████╗ ███████╗████████╗███████╗██████╗
██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗
█████╗  ███████║███████╗   ██║   █████╗  ██████╔╝
██╔══╝  ██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗
██║     ██║  ██║███████║   ██║   ███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
`

	// Create the UI
	ui := strings.Join([]string{
		"\n\n",
		logo,
		// "\n\n",
		// titleStyle.Render("FASTER"),
		"\n\n",
		subtitleStyle.Render("Press any key to continue..."),
		"\n",
	}, "\n")

	return ui
}
