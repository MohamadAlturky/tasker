package main

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	now := time.Now()
	task := NewTask(todo, "Test Task", "Test Description", now)

	if task.status != todo {
		t.Errorf("Expected status %v, got %v", todo, task.status)
	}
	if task.title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.title)
	}
	if task.description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", task.description)
	}
	if !task.dueDate.Equal(now) {
		t.Errorf("Expected due date %v, got %v", now, task.dueDate)
	}
}

func TestTaskNext(t *testing.T) {
	tests := []struct {
		name     string
		initial  status
		expected status
	}{
		{"todo to inProgress", todo, inProgress},
		{"inProgress to done", inProgress, done},
		{"done to todo", done, todo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{status: tt.initial}
			task.Next()
			if task.status != tt.expected {
				t.Errorf("Expected status %v after Next(), got %v", tt.expected, task.status)
			}
		})
	}
}

func TestTaskFilterValue(t *testing.T) {
	task := Task{title: "Test Task"}
	if task.FilterValue() != "Test Task" {
		t.Errorf("Expected FilterValue() to return 'Test Task', got '%s'", task.FilterValue())
	}
}

func TestTaskTitleAndDescription(t *testing.T) {
	task := Task{
		title:       "Test Title",
		description: "Test Description",
	}

	if task.Title() != "Test Title" {
		t.Errorf("Expected Title() to return 'Test Title', got '%s'", task.Title())
	}

	if task.Description() != "Test Description" {
		t.Errorf("Expected Description() to return 'Test Description', got '%s'", task.Description())
	}
}
