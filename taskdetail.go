package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TaskDetail displays detailed information about a task
type TaskDetail struct {
	task    Task
	width   int
	height  int
	visible bool
}

// NewTaskDetail creates a new task detail view
func NewTaskDetail(task Task) TaskDetail {
	return TaskDetail{
		task:    task,
		width:   80, // default width
		height:  24, // default height
		visible: true,
	}
}

func (td TaskDetail) Init() tea.Cmd {
	return nil
}

func (td TaskDetail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		td.width = msg.Width
		td.height = msg.Height
		return td, nil
	case tea.KeyMsg:
		// Any key press returns to the board view
		return board, nil
	}
	return td, nil
}

func (td TaskDetail) View() string {
	if !td.visible {
		return ""
	}

	// Calculate container width based on terminal width
	containerWidth := td.width * 3 / 4
	if containerWidth < 60 {
		containerWidth = 60
	}
	if containerWidth > td.width-20 {
		containerWidth = td.width - 20
	}

	// Define styles
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2).
		Width(containerWidth)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginBottom(1).
		Italic(true).
		Width(containerWidth - 10)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("110")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(containerWidth - 10).
		MarginTop(1).
		MarginBottom(1)

	// Add status color based on task status
	var statusColor string
	var statusText string
	switch td.task.status {
	case todo:
		statusText = "To Do"
		statusColor = "62" // Blue
	case inProgress:
		statusText = "In Progress"
		statusColor = "208" // Orange
	case done:
		statusText = "Done"
		statusColor = "34" // Green
	default:
		statusText = "Unknown"
		statusColor = "240" // Gray
	}

	// Create the content
	title := titleStyle.Render(td.task.title)
	divider := dividerStyle.Render("─────────────────────────────────────")

	statusLabel := labelStyle.Render("Status:")
	statusValue := valueStyle.Copy().Foreground(lipgloss.Color(statusColor)).Bold(true).Render(statusText)
	status := lipgloss.JoinHorizontal(lipgloss.Left, statusLabel, "  ", statusValue)

	idLabel := labelStyle.Render("ID:")
	idValue := valueStyle.Render(fmt.Sprintf("%d", td.task.id))
	id := lipgloss.JoinHorizontal(lipgloss.Left, idLabel, "  ", idValue)

	descriptionLabel := labelStyle.Render("Description:")
	descriptionText := td.task.description
	if descriptionText == "" {
		descriptionText = "No description provided"
	}

	description := lipgloss.JoinVertical(
		lipgloss.Left,
		descriptionLabel,
		valueStyle.Render(descriptionText),
	)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Padding(1, 0, 0, 0).
		Render("Press any key to go back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		divider,
		status,
		id,
		"",
		description,
		divider,
		instructions,
	)

	// Render the container with content
	renderedContent := containerStyle.Render(content)

	// Center the content
	return lipgloss.Place(
		td.width,
		td.height,
		lipgloss.Center,
		lipgloss.Center,
		renderedContent,
	)
}
