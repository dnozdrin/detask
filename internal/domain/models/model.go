package models

import "time"

// Model represents the default fields for persisted structures
type Model struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// Board represents a board (project)
type Board struct {
	Model
	Name        string `json:"name" validate:"max=500,min=1"`
	Description string `json:"description" validate:"max=1000"`
}

// Column represents a column (status)
type Column struct {
	Model
	Name     string  `json:"name" validate:"max=255,min=1"`
	BoardID  uint    `json:"board"`
	Position float64 `json:"position"`
}

// Task represents a task
type Task struct {
	Model
	Name        string  `json:"name,omitempty" validate:"max=500,min=1"`
	Description string  `json:"description,omitempty" validate:"max=5000"`
	ColumnID    uint    `json:"column,omitempty"`
	Position    float64 `json:"column,omitempty"`
}

// Comment represents a comment to a task
type Comment struct {
	Model
	Text   string `json:"text" validate:"max=5000,min=1"`
	TaskID uint   `json:"task"`
}
