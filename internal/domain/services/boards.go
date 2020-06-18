// todo: consider patch implementation
// todo: consider adding filter by constraints
package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

type BoardStorage interface {
	Save(board *m.Board) (*m.Board, error)
	FindById(ID uint) (*m.Board, error)
	FindAll() ([]*m.Board, error)
	Update(board *m.Board) (*m.Board, error)
	Delete(ID uint) error
}

type BoardService struct {
	boardStorage BoardStorage
	validator    v.Validator
}

// NewBoardService is a board service constructor
func NewBoardService(boardStorage BoardStorage, validator v.Validator) *BoardService {
	return &BoardService{
		boardStorage: boardStorage,
		validator:    validator,
	}
}

// Create will create a new board with the provided payload
// and a default column
func (b *BoardService) Create(board *m.Board) (*m.Board, v.Result) {
	var result v.Result
	if result = b.validator.Validate(*board); !result.IsValid() {
		return nil, result
	}
	// todo: create a default column here

	board, result.Error = b.boardStorage.Save(board)

	return board, result
}

// FindAll will return all not deleted boards and an error in case
// it occurred while fetching records from the storage
func (b *BoardService) FindAll() ([]*m.Board, error) {
	return b.boardStorage.FindAll()
}

// FindOneById will return a pointer to the board requested by id and
// an error in case it occurred while fetching the record from the storage
func (b *BoardService) FindOneById(ID uint) (*m.Board, error) {
	return b.boardStorage.FindById(ID)
}

// Update will update the board record
func (b *BoardService) Update(board *m.Board) (*m.Board, v.Result) {
	var result v.Result
	if result = b.validator.Validate(*board); !result.IsValid() {
		return nil, result
	}

	board, result.Error = b.boardStorage.Update(board)

	return board, result
}

// Delete will delete a record with the given ID
func (b *BoardService) Delete(ID uint) error {
	return b.boardStorage.Delete(ID)
}
