package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1
const Move = -2

type column struct {
	focus  bool
	status status
	list   list.Model
	height int
	width  int
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status status) column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus: focus, status: status, list: defaultList}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewForm(task.title, task.description, task.dueDate)
				f.index = c.list.Index()
				f.col = c
				f.taskID = task.GetID()
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) == 0 {
		return nil
	}

	// Get the current task
	task := c.list.SelectedItem().(Task)

	// Delete from database
	err := database.DeleteTask(task.GetID())
	if err != nil {
		// Return an error message
		return func() tea.Msg {
			return tea.Printf("Error deleting task: %v", err)
		}
	}

	// Remove from the list
	c.list.RemoveItem(c.list.Index())

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i == Move {
		// Just insert the task at the end of the list without database operations
		// since it's being moved from another column and already in the database
		return c.list.InsertItem(APPEND, t)
	}
	if i != APPEND {
		// Update existing task in database
		err := database.UpdateTask(t)
		if err != nil {
			return func() tea.Msg {
				return tea.Printf("Error updating task: %v", err)
			}
		}
		return c.list.SetItem(i, t)
	}

	// Create new task in database
	// The issue was that database.CreateTask doesn't update the original task
	// because it receives a copy, not a pointer
	// So we need to get the ID separately and set it on our task
	err := database.CreateTask(t)
	if err != nil {
		return func() tea.Msg {
			return tea.Printf("Error creating task: %v", err)
		}
	}

	// Get the latest tasks to find our newly created task
	tasks, err := database.GetAllTasks()
	if err != nil {
		return func() tea.Msg {
			return tea.Printf("Error retrieving tasks: %v", err)
		}
	}

	// Find the newly created task (assuming it's the last one with matching title)
	for i := len(tasks) - 1; i >= 0; i-- {
		if tasks[i].title == t.title && tasks[i].status == t.status {
			t.SetID(tasks[i].GetID())
			break
		}
	}

	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width, height int) {
	c.width = width / margin
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

// moveMsg struct definition will be moved to messages.go
// keep the MoveToNext function that uses it

func (c *column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}

	// Update task status
	nextStatus := c.status.getNext()
	task.status = nextStatus

	// Remove item from the current list
	c.list.RemoveItem(c.list.Index())

	// Refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	// Return the moveMsg to add the task to its new column
	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}
