package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         column
	index       int
	taskID      int
	confirming  bool
	width       int
	height      int
	isEdit      bool
}

func newDefaultForm() *Form {
	return NewForm("task name", "", false)
}

func NewForm(title, description string, isEdit bool) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
		taskID:      -1,
		width:       80, // default width, will be updated when window size is received
		height:      24, // default height, will be updated when window size is received
		isEdit:      isEdit,
	}
	form.title.Placeholder = title
	form.description.Placeholder = description

	// Set actual value for title and description when editing
	if isEdit {
		form.title.SetValue(title)
		form.description.SetValue(description)
	}

	form.title.Focus()
	return &form
}

func (f Form) CreateTask() Task {
	task := Task{
		status:      f.col.status,
		title:       f.title.Value(),
		description: f.description.Value(),
	}
	if f.taskID != -1 {
		task.SetID(f.taskID)
	}
	return task
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		return f, nil
	case column:
		f.col = msg
		f.col.list.Index()
	case tea.KeyMsg:
		if f.confirming {
			switch msg.String() {
			case "y", "Y", "enter":
				return board.Update(f)
			case "n", "N", "esc":
				f.confirming = false
				return f, nil
			default:
				return f, nil
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			if msg.String() == "ctrl+c" {
				return f, tea.Quit
			}
		case key.Matches(msg, keys.Back):
			if f.description.Focused() {
				f.confirming = true
				return f, nil
			}
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			if f.description.Focused() {
				f.confirming = true
				return f, nil
			}
			return board.Update(f)
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	// Calculate the form width based on terminal width
	formWidth := f.width - 4 // Subtract padding and borders
	if formWidth < 20 {      // Minimum width
		formWidth = 20
	}

	// Calculate textarea width
	textareaWidth := formWidth - 4 // Subtract padding

	// Update textarea width
	f.description.SetWidth(textareaWidth)

	formStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(formWidth)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1).
		MarginBottom(1)

	descriptionLabel := "Description"
	if f.description.Focused() && !f.confirming {
		descriptionLabel = "Description (press ESC to save or Enter for confirmation)"
	}

	formTitle := "Create a new task"
	if f.isEdit {
		formTitle = "Edit task"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(formTitle),
		labelStyle.Render("Task Title"),
		f.title.View(),
		labelStyle.Render(descriptionLabel),
		f.description.View(),
		f.help.View(keys),
	)

	renderedContent := formStyle.Render(content)

	// If we're in confirmation mode, create a popup
	if f.confirming {
		// Define styles for the popup
		popupWidth := formWidth / 2
		if popupWidth < 30 {
			popupWidth = 30
		}

		popupStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 3).
			Foreground(lipgloss.Color("15")).
			Width(popupWidth).
			Align(lipgloss.Center)

		// Create the popup content
		popupContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"Save changes?",
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(y/n or enter/esc)"),
		)

		// Render the popup
		popup := popupStyle.Render(popupContent)

		// Center the popup
		centered := lipgloss.Place(
			f.width,
			f.height,
			lipgloss.Center,
			lipgloss.Center,
			popup,
		)

		return centered
	}

	return renderedContent
}
