package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite", "./tasks.db")
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
	_, err := d.db.Exec(`
		INSERT INTO tasks (title, description, due_date, status)
		VALUES (?, ?, ?, ?)
	`, task.title, task.description, task.dueDate, task.status)
	return err
}

func (d *Database) UpdateTask(task Task) error {
	_, err := d.db.Exec(`
		UPDATE tasks
		SET title = ?,
			description = ?,
			due_date = ?,
			status = ?
		WHERE id = ?
	`, task.title, task.description, task.dueDate, task.status, task.id)
	return err
}

func (d *Database) DeleteTask(id int) error {
	_, err := d.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (d *Database) GetAllTasks() ([]Task, error) {
	rows, err := d.db.Query(`
		SELECT id, title, description, due_date, status
		FROM tasks
		ORDER BY status, due_date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		err := rows.Scan(&task.id, &task.title, &task.description, &dueDate, &task.status)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			task.dueDate = dueDate.Time
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
