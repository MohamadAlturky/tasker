package main

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *Database {
	// Create a temporary database file
	dbPath := "test_tasks.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			due_date DATETIME,
			status INTEGER NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return &Database{db: db}
}

func cleanupTestDB(t *testing.T, db *Database) {
	// Close the database connection
	if db != nil {
		db.Close()
	}
	// Remove the test database file
	os.Remove("test_tasks.db")
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	now := time.Now()
	task := NewTask(todo, "Test Task", "Test Description", now)

	err := db.CreateTask(task)
	if err != nil {
		t.Errorf("Failed to create task: %v", err)
	}

	// Verify the task was created
	tasks, err := db.GetAllTasks()
	if err != nil {
		t.Errorf("Failed to get tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].title != task.title {
		t.Errorf("Expected task title '%s', got '%s'", task.title, tasks[0].title)
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create initial task
	now := time.Now()
	task := NewTask(todo, "Test Task", "Test Description", now)
	err := db.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Get the created task to get its ID
	tasks, err := db.GetAllTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	task = tasks[0]

	// Update task status
	task.status = inProgress
	err = db.UpdateTask(task)
	if err != nil {
		t.Errorf("Failed to update task: %v", err)
	}

	// Verify the update
	tasks, err = db.GetAllTasks()
	if err != nil {
		t.Errorf("Failed to get tasks: %v", err)
	}
	if tasks[0].status != inProgress {
		t.Errorf("Expected status %d, got %d", inProgress, tasks[0].status)
	}
}

func TestGetAllTasks(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create test tasks
	now := time.Now()
	tasks := []Task{
		NewTask(todo, "Task 1", "Description 1", now),
		NewTask(inProgress, "Task 2", "Description 2", now),
		NewTask(done, "Task 3", "Description 3", now),
	}

	for _, task := range tasks {
		err := db.CreateTask(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
	}

	// Get all tasks
	allTasks, err := db.GetAllTasks()
	if err != nil {
		t.Errorf("Failed to get tasks: %v", err)
	}
	if len(allTasks) != len(tasks) {
		t.Errorf("Expected %d tasks, got %d", len(tasks), len(allTasks))
	}
}
