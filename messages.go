package main

// Message types for communication between components

// moveMsg is sent when a task is moved from one column to another
type moveMsg struct {
	Task
}
