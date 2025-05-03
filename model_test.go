package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	if board == nil {
		t.Fatal("Expected NewBoard() to return a non-nil Board")
	}
	if board.focused != todo {
		t.Errorf("Expected initial focus to be todo, got %v", board.focused)
	}
	if !board.help.ShowAll {
		t.Error("Expected help.ShowAll to be true")
	}
}

func TestBoardInit(t *testing.T) {
	board := NewBoard()
	cmd := board.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil command")
	}
}

func TestBoardUpdateWindowSize(t *testing.T) {
	board := NewBoard()
	// Initialize columns
	board.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	_, _ = board.Update(msg)
	if !board.loaded {
		t.Error("Expected board to be marked as loaded after WindowSizeMsg")
	}
	if board.help.Width != msg.Width-margin {
		t.Errorf("Expected help width to be %d, got %d", msg.Width-margin, board.help.Width)
	}
}

func TestBoardUpdateQuit(t *testing.T) {
	board := NewBoard()
	// Initialize columns
	board.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := board.Update(msg)
	if cmd == nil {
		t.Error("Expected Update with quit key to return a non-nil command")
	}
	if !board.quitting {
		t.Error("Expected board to be marked as quitting after quit key")
	}
}

func TestBoardView(t *testing.T) {
	board := NewBoard()
	// Initialize columns
	board.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	// Test quitting state
	board.quitting = true
	if board.View() != "" {
		t.Error("Expected empty string for quitting state")
	}

	// Test loading state
	board.quitting = false
	board.loaded = false
	if board.View() != "loading..." {
		t.Error("Expected 'loading...' for loading state")
	}

	// Test normal state
	board.loaded = true
	view := board.View()
	if view == "" {
		t.Error("Expected non-empty view for loaded state")
	}
}
