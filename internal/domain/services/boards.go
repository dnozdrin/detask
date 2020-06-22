// todo: consider adding filter by constraints
// todo: add error wrappers
package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// BoardStorage represents an interface for interaction with boards DAO
type BoardStorage interface {
	// SaveWithDefaultColumn should persist the provided board and create a default
	// column for it
	SaveWithDefaultColumn(board *m.Board) (*m.Board, error)
	// FindOneById should return a board with the provided ID
	FindOneById(ID uint) (*m.Board, error)
	// Find should return a slice of boards pointers sorted by name, that meet the
	// provided demand
	Find() ([]*m.Board, error)
	// Update should update all board fields by the provided data
	Update(board *m.Board) (*m.Board, error)
	// Delete should set current deletion time to a board with the provided ID
	// and to all dependant records
	Delete(ID uint) error
}

// BoardService is an interactor for work with boards
type BoardService struct {
	validator     v.Validator
	boardStorage  BoardStorage
	columnStorage ColumnStorage
}

// NewBoardService is a board service constructor
func NewBoardService(validator v.Validator, boardStorage BoardStorage) *BoardService {
	return &BoardService{
		validator:    validator,
		boardStorage: boardStorage,
	}
}

// Create will create a new board with the provided payload
// and a default column
func (b *BoardService) Create(board *m.Board) (*m.Board, error) {
	if err := b.validator.Validate(*board); err != nil {
		return nil, err
	}

	return b.boardStorage.SaveWithDefaultColumn(board)
}

// Find will return all not deleted boards and an error in case
// it occurred while fetching records from the storage
func (b *BoardService) Find() ([]*m.Board, error) {
	return b.boardStorage.Find()
}

// FindOneById will return a pointer to the board requested by id and
// an error in case it occurred while fetching the record from the storage
func (b *BoardService) FindOneById(ID uint) (*m.Board, error) {
	return b.boardStorage.FindOneById(ID)
}

// Update will update the board record
func (b *BoardService) Update(board *m.Board) (*m.Board, error) {
	if err := b.validator.Validate(*board); err != nil {
		return nil, err
	}

	return b.boardStorage.Update(board)
}

// Delete will mark a record with the given ID as deleted as well as all
// the dependant records
func (b *BoardService) Delete(ID uint) error {
	return b.boardStorage.Delete(ID)
}
