package main

import "time"

type Task struct {
	id          int
	status      status
	title       string
	description string
	dueDate     time.Time
}

func NewTask(status status, title, description string, dueDate time.Time) Task {
	return Task{
		status:      status,
		title:       title,
		description: description,
		dueDate:     dueDate,
	}
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
