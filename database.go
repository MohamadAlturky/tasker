package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	// Get the current executable's directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	appDir := filepath.Dir(exePath)

	// Use the app directory for the database file
	dbPath := filepath.Join(appDir, "tasks.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
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
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateTask(task Task) error {
	// Validate the task before inserting
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %v", err)
	}

	result, err := d.db.Exec(`
		INSERT INTO tasks (title, description, due_date, status)
		VALUES (?, ?, ?, ?)
	`, task.title, task.description, task.dueDate, task.status)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	// Get the ID of the newly created task
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get task ID: %v", err)
	}
	task.SetID(int(id))

	return nil
}

func (d *Database) UpdateTask(task Task) error {
	// Validate the task before updating
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %v", err)
	}

	// Check if task exists
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", task.GetID()).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if task exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("task with id %d not found", task.GetID())
	}

	_, err = d.db.Exec(`
		UPDATE tasks
		SET title = ?,
			description = ?,
			due_date = ?,
			status = ?
		WHERE id = ?
	`, task.title, task.description, task.dueDate, task.status, task.GetID())
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

func (d *Database) DeleteTask(id int) error {
	// First verify the task exists
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if task exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("task with id %d not found", id)
	}

	// Delete the task
	result, err := d.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	// Verify that a row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task was deleted (id: %d)", id)
	}

	return nil
}

func (d *Database) GetAllTasks() ([]Task, error) {
	rows, err := d.db.Query(`
		SELECT id, title, description, due_date, status
		FROM tasks
		ORDER BY status, due_date
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		err := rows.Scan(&task.id, &task.title, &task.description, &dueDate, &task.status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		if dueDate.Valid {
			task.dueDate = dueDate.Time
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	return tasks, nil
}
