package models

import (
	"time"

	"github.com/pkg/errors"
)

type Model struct {
	ID        uint       `json:"id" db:"id"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

var (
	ErrRecordNotFound     = errors.New("record was not found")
	ErrRecordAlreadyExist = errors.New("this records already exists")
)

type Board struct {
	Model
	Name        string `db:"name" json:"name" validate:"max=500,min=1"`
	Description string `db:"description" json:"description" validate:"max=1000"`
}

type Column struct {
	Model
	Name     string `json:"name" validate:"max=255,min=1"`
	BoardID  uint   `db:"board" json:"board"`
	Position uint   `db:"position" json:"position"`
}

type Task struct {
	Model
	Name        string  `db:"name" json:"name,omitempty" validate:"max=500,min=1"`
	Description string  `db:"description" json:"description,omitempty" validate:"max=5000"`
	ColumnID    uint    `db:"column" json:"column,omitempty"`
	Position    float64 `db:"position" json:"column,omitempty"`
}

type Comment struct {
	Model
	Text   string `db:"text" json:"text" validate:"max=5000,min=1"`
	TaskID uint   `db:"task" json:"task"`
}
