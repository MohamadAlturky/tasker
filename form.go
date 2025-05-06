package main

import (
	"time"

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
	dueDate     textinput.Model
	col         column
	index       int
	taskID      int
	confirming  bool
}

func newDefaultForm() *Form {
	return NewForm("task name", "", time.Now())
}

func NewForm(title, description string, dueDate time.Time) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
		dueDate:     textinput.New(),
		taskID:      -1,
	}
	form.title.Placeholder = title
	form.description.Placeholder = description
	form.description.SetValue(description)
	form.dueDate.Placeholder = "YYYY-MM-DD"
	form.dueDate.SetValue(dueDate.Format("2006-01-02"))
	form.title.Focus()
	return &form
}

func (f Form) CreateTask() Task {
	dueDate, _ := time.Parse("2006-01-02", f.dueDate.Value())
	task := Task{
		status:      f.col.status,
		title:       f.title.Value(),
		description: f.description.Value(),
		dueDate:     dueDate,
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
			return f, tea.Quit
		case key.Matches(msg, keys.Back):
			if f.description.Focused() {
				return board.Update(f)
			}
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.dueDate.Focus()
				return f, textinput.Blink
			}
			if f.dueDate.Focused() {
				f.dueDate.Blur()
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
	if f.dueDate.Focused() {
		f.dueDate, cmd = f.dueDate.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	formStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(50)

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

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Create a new task"),
		labelStyle.Render("Task Title"),
		f.title.View(),
		labelStyle.Render("Due Date (YYYY-MM-DD)"),
		f.dueDate.View(),
		labelStyle.Render(descriptionLabel),
		f.description.View(),
		f.help.View(keys),
	)

	renderedContent := formStyle.Render(content)

	// If we're in confirmation mode, create a popup
	if f.confirming {
		// Define styles for the popup
		popupStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 3).
			Foreground(lipgloss.Color("15")).
			Width(30).
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
			formStyle.GetWidth(),
			20,
			lipgloss.Center,
			lipgloss.Center,
			popup,
		)

		return centered
	}

	return renderedContent
}
