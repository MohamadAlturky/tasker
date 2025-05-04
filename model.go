package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help     help.Model
	loaded   bool
	focused  status
	cols     []column
	quitting bool
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

// MoveTask moves a task from one status to another and updates the database
func (m *Board) MoveTask(task Task, newStatus status) tea.Cmd {
	// Update task status
	task.status = newStatus

	// Update in database
	err := database.UpdateTask(task)
	if err != nil {
		return func() tea.Msg {
			return tea.Printf("Error updating task: %v", err)
		}
	}

	// Add to destination column using the Move constant
	return m.cols[newStatus].Set(Move, task)
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case Form:
		task := msg.CreateTask()
		// err := database.CreateTask(task)
		// if err != nil {
		// 	// Log the error and show it to the user
		// 	return m, func() tea.Msg {
		// 		return tea.Printf("Error creating task: %v", err)
		// 	}
		// }
		return m, m.cols[m.focused].Set(msg.index, task)
	case moveMsg:
		// Task was already removed from source list in column.MoveToNext
		// Just need to move it to the destination column
		return m, m.MoveTask(msg.Task, msg.Task.status)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
	} else {
		return res, cmd
	}
	return m, cmd
}

func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}
