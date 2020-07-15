package services

import (
	"database/sql"
	m "github.com/dnozdrin/detask/internal/domain/models"
)

// BoardStorage represents an interface for interaction with boards DAO
type BoardStorage interface {
	// Save should persist the provided board and create a default
	// column for it
	Save(*m.Board) (*m.Board, error)
	// FindOneById should return a board with the provided ID
	FindOneById(uint) (*m.Board, error)
	// Find should return a slice of boards pointers sorted by name, that meet the
	// provided demand
	Find() ([]*m.Board, error)
	// Update should update all board fields by the provided data
	Update(*m.Board) (*m.Board, error)
	// Delete should delete a board with the provided ID as well as all dependant records
	Delete(uint) error
	// WithTx should return the boardStorage that will use the provided transaction
	WithTx(*sql.Tx) BoardStorage
}

// ColumnStorage represents an interface for interaction with columns DAO
type ColumnStorage interface {
	// Save will persist the provided column
	Save(*m.Column) (*m.Column, error)
	// FindOneById should return a column with the provided ID
	FindOneById(uint) (*m.Column, error)
	// Find should return a slice of columns pointers sorted by position, that meet the
	// provided demand
	Find(ColumnDemand) ([]*m.Column, error)
	// Update should update all column fields by the provided data
	Update(*m.Column) (*m.Column, error)
	// Delete should delete a column with the provided ID as well as all dependant records
	Delete(uint) error
	// WithTx should return the columnStorage that will use the provided transaction
	WithTx(*sql.Tx) ColumnStorage
	// CountColumnsByBoard should count columns that are related to provided board ID
	CountColumnsByBoard(uint) (int, error)
	// FindColumnToTheLeft should find a ID of a column that is to the left of the current
	// and is related to the same board
	FindColumnToTheLeft(uint) (uint, error)
	// FindColumnToTheRight should find a ID of a column that is to the right of the current
	// and is related to the same board
	FindColumnToTheRight(uint) (uint, error)
}

// TaskStorage represents an interface for interaction with tasks DAO
type TaskStorage interface {
	// Save will persist the provided task
	Save(*m.Task) (*m.Task, error)
	// FindOneById should return a task with the provided ID
	FindOneById(uint) (*m.Task, error)
	// Find should return a slice of boards pointers sorted by name, that meet the
	// provided demand
	Find(TaskDemand) ([]*m.Task, error)
	// Update should update the name and the description of the task
	Update(*m.Task) (*m.Task, error)
	// Delete should delete a task with the provided ID as well as all dependant records
	Delete(uint) error
	// WithTx should return the taskStorage that will use the provided transaction
	WithTx(*sql.Tx) TaskStorage
	// MoveToColumn should move all task from one column to another
	MoveToColumn(from, to uint) error
}

// CommentStorage represents an interface for interaction with comments DAO
type CommentStorage interface {
	// Save will persist the provided comment
	Save(*m.Comment) (*m.Comment, error)
	// FindOneById should return a comment with the provided ID
	FindOneById(uint) (*m.Comment, error)
	// Find should return a slice of comments pointers sorted by creation date
	// (from newest to oldest), that meet the provided demand
	Find(CommentDemand) ([]*m.Comment, error)
	// Update should update the comment text
	Update(*m.Comment) (*m.Comment, error)
	// Delete should delete a comment with the provided ID as well as all dependant records
	Delete(uint) error
}

// TxBeginner provides a method for starting database transactions
type TxBeginner interface {
	Begin() (*sql.Tx, error)
}
