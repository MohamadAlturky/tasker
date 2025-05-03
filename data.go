package main

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
		// Handle error appropriately
		return
	}

	// Add tasks to appropriate columns
	for _, task := range tasks {
		b.cols[task.status].list.SetItems(append(b.cols[task.status].list.Items(), task))
	}
}
