package main

import (
	"fmt"
	"time"
)

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

// SetID sets the task's ID
func (t *Task) SetID(id int) {
	t.id = id
}

// GetID returns the task's ID
func (t Task) GetID() int {
	return t.id
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// Validate checks if the task has all required fields
func (t Task) Validate() error {
	if t.title == "" {
		return fmt.Errorf("task title cannot be empty")
	}
	if t.status < todo || t.status > done {
		return fmt.Errorf("invalid task status")
	}
	return nil
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
