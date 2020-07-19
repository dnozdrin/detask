package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// BoardService is an interactor for work with boards
type BoardService struct {
	validator     v.Validator
	boardStorage  BoardStorage
	columnStorage ColumnStorage
	txBeginner    TxBeginner
}

// NewBoardService is a board service constructor
func NewBoardService(
	validator v.Validator,
	boardStorage BoardStorage,
	columnStorage ColumnStorage,
	txBeginner TxBeginner,
) *BoardService {
	return &BoardService{
		validator:     validator,
		boardStorage:  boardStorage,
		columnStorage: columnStorage,
		txBeginner:    txBeginner,
	}
}

// Create will create a new board with the provided payload
// and a default column
func (b *BoardService) Create(board *m.Board) (*m.Board, error) {
	if err := b.validator.Validate(*board); err != nil {
		return nil, err
	}

	tx, err := b.txBeginner.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	boardStorage := b.boardStorage.WithTx(tx)
	board, err = boardStorage.Save(board)
	if err != nil {
		return nil, err
	}

	column := &m.Column{
		Name:     "Default",
		BoardID:  board.ID,
		Position: DefaultColPos,
	}
	columnStorage := b.columnStorage.WithTx(tx)
	if _, err = columnStorage.Save(column); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return board, nil
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
