package main

import (
	"math/rand"
	"strings"
	"time"

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

// Returns a random color from a predefined set of vibrant colors
func getRandomColor() lipgloss.Color {
	rand.Seed(time.Now().UnixNano())
	colors := []lipgloss.Color{
		"#FF6B6B", // Red
		"#4ECDC4", // Teal
		"#FFE66D", // Yellow
		"#6BFF84", // Green
		"#FF85EB", // Pink
		"#85C1FF", // Blue
		"#C385FF", // Purple
		"#FF9B85", // Orange
	}
	return colors[rand.Intn(len(colors))]
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
		Foreground(lipgloss.Color(getRandomColor())).
		Align(lipgloss.Center)

	// Create the logo with random color
	logoStyle := lipgloss.NewStyle().
		Foreground(getRandomColor()).
		Bold(true)

	logoText := `
███████╗ █████╗ ███████╗████████╗███████╗██████╗
██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗
█████╗  ███████║███████╗   ██║   █████╗  ██████╔╝
██╔══╝  ██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗
██║     ██║  ██║███████║   ██║   ███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
`
	logo := logoStyle.Render(logoText)

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
