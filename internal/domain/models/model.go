package models

import "time"

// Model represents the default fields for persisted structures
type Model struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Board represents a board (project)
type Board struct {
	Model
	Name        string `json:"name" validate:"required,max=500,min=1"`
	Description string `json:"description" validate:"max=1000"`
}

// Column represents a column (status)
type Column struct {
	Model
	Name     string  `json:"name" validate:"required,max=255,min=1"`
	BoardID  uint    `json:"board" validate:"required,numeric"`
	Position float64 `json:"position" validate:"required,numeric"`
}

// Task represents a task
type Task struct {
	Model
	Name        string  `json:"name" validate:"required,max=500,min=1"`
	Description string  `json:"description" validate:"max=5000"`
	ColumnID    uint    `json:"column" validate:"required,numeric"`
	Position    float64 `json:"position" validate:"required,numeric"`
}

// Comment represents a comment to a task
type Comment struct {
	Model
	Text   string `json:"text" validate:"required,max=5000,min=1"`
	TaskID uint   `json:"task" validate:"required,numeric"`
}
