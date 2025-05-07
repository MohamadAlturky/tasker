package main

import tea "github.com/charmbracelet/bubbletea"

// Provides the data to fill the kanban board from the database

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	// Set column titles
	b.cols[todo].list.Title = "To Do"
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[done].list.Title = "Done"

	// Load tasks from database
	tasks, err := database.GetAllTasks()
	if err != nil {
		// Log the error and show it to the user
		tea.Printf("Error loading tasks: %v", err)
		return
	}

	// Add tasks to appropriate columns
	for _, task := range tasks {
		// Skip archived tasks or convert them to done status
		if task.status > done {
			task.status = done
		}
		b.cols[task.status].list.SetItems(append(b.cols[task.status].list.Items(), task))
	}
}
