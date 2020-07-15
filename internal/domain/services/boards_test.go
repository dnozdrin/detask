// +build unit

package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardService(t *testing.T) {
	boardStorage := new(MockedBoardStorage)
	columnStorage := new(MockedColumnStorage)
	validation := new(MockedValidation)
	txBeginner := new(MockedTxBeginner)
	boardService := NewBoardService(validation, boardStorage, columnStorage, txBeginner)

	assert.Equal(t, validation, boardService.validator)
	assert.Equal(t, boardStorage, boardService.boardStorage)
	assert.Equal(t, columnStorage, boardService.columnStorage)
	assert.Equal(t, txBeginner, boardService.txBeginner)
}

func TestBoardService_Create(t *testing.T) {
	var boardIn = &m.Board{Name: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		savedBoard := &m.Board{Model: m.Model{ID: 123}, Name: "dummy"}
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(savedBoard, nil)
		boardStorage.On("WithTx", tx).Return(boardStorage)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		column := &m.Column{Name: "Default", Position: DefaultColPos, BoardID: savedBoard.ID}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", column).Return(column, nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		boardService := &BoardService{
			validator:     validation,
			boardStorage:  boardStorage,
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}

		resultBoard, err := boardService.Create(boardIn)

		assert.NotNil(t, resultBoard)
		assert.Nil(t, err)
		assert.Equal(t, savedBoard, resultBoard)
	})
	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{validator: validation}
		boardOut, resultOut := boardService.Create(boardIn)

		assert.Equal(t, validationErr, resultOut)
		assert.Empty(t, boardOut)
	})
	t.Run("board_save_error", func(t *testing.T) {
		var validationErr *v.Errors
		dbErr := errors.New("simple error")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectRollback()
		tx, _ := db.Begin()

		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(&m.Board{}, dbErr)
		boardStorage.On("WithTx", tx).Return(boardStorage)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		boardService := &BoardService{
			validator:    validation,
			boardStorage: boardStorage,
			txBeginner:   txBeginner,
		}
		boardOut, err := boardService.Create(boardIn)

		assert.Empty(t, boardOut)
		assert.Equal(t, dbErr, err)
	})
	t.Run("column_save_error", func(t *testing.T) {
		dbErr := errors.New("simple error")
		var validationErr *v.Errors

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectRollback()
		tx, _ := db.Begin()

		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(&m.Board{Model: m.Model{ID: 123}}, nil)
		boardStorage.On("WithTx", tx).Return(boardStorage)

		column := &m.Column{Name: "Default", Position: DefaultColPos, BoardID: 123}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", column).Return(&m.Column{}, dbErr)
		columnStorage.On("WithTx", tx).Return(columnStorage)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		boardService := &BoardService{
			validator:     validation,
			boardStorage:  boardStorage,
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}
		boardOut, err := boardService.Create(boardIn)

		assert.Empty(t, boardOut)
		assert.Equal(t, dbErr, err)
	})
	t.Run("transaction_begin_error", func(t *testing.T) {
		txErr := errors.New("tx error")
		var validationErr *v.Errors

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		tx, _ := db.Begin()

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, txErr)

		boardService := &BoardService{
			validator:  validation,
			txBeginner: txBeginner,
		}
		boardOut, err := boardService.Create(boardIn)

		assert.Empty(t, boardOut)
		assert.Equal(t, txErr, err)
	})

	t.Run("transaction_commit_error", func(t *testing.T) {
		txErr := errors.New("tx error")
		var validationErr *v.Errors

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit().WillReturnError(txErr)
		tx, _ := db.Begin()

		savedBoard := &m.Board{Model: m.Model{ID: 123}, Name: "dummy"}
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(savedBoard, nil)
		boardStorage.On("WithTx", tx).Return(boardStorage)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		column := &m.Column{Name: "Default", Position: DefaultColPos, BoardID: savedBoard.ID}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", column).Return(column, nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		boardService := &BoardService{
			validator:     validation,
			boardStorage:  boardStorage,
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}

		resultBoard, err := boardService.Create(boardIn)
		assert.Empty(t, resultBoard)
		assert.Equal(t, txErr, err)
	})
}

func TestBoardService_FindOneById(t *testing.T) {
	const dummyID = 1234
	boardIn := &m.Board{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindOneById", mock.Anything).Return(boardIn, nil)
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, boardIn, boardOut)
	})

	t.Run("not_found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindOneById", mock.Anything).Return(boardIn, errors.New(""))
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, boardIn, boardOut)
	})
}

func TestBoardService_Find(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		boardsIn := []*m.Board{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Find").Return(boardsIn, nil)
		boardService := &BoardService{boardStorage: boardStorage}
		boardsOut, err := boardService.Find()
		assert.Nil(t, err)
		assert.Equal(t, boardsIn, boardsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Find", mock.Anything).Return([]*m.Board{}, errors.New(mock.Anything))
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.Find()
		assert.Error(t, err)
		assert.Empty(t, boardOut)
	})
}

func TestBoardService_Update(t *testing.T) {
	var boardIn = &m.Board{Name: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Update", boardIn).Return(boardIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{
			boardStorage: boardStorage,
			validator:    validation,
		}
		boardOut, err := boardService.Update(boardIn)

		assert.NotNil(t, boardOut)
		assert.Nil(t, err)
	})

	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{validator: validation}
		boardOut, resultOut := boardService.Update(boardIn)

		assert.Equal(t, validationErr, resultOut)
		assert.Empty(t, boardOut)
	})

	t.Run("database_error", func(t *testing.T) {
		dbErr := errors.New("simple error")
		var validationErr *v.Errors
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Update", boardIn).Return(&m.Board{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{
			boardStorage: boardStorage,
			validator:    validation,
		}
		boardOut, err := boardService.Update(boardIn)

		assert.Empty(t, boardOut)
		assert.Equal(t, err, dbErr)
	})
}

func TestBoardService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Delete", mock.Anything).Return(nil)
		boardService := &BoardService{boardStorage: boardStorage}
		err := boardService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Delete", mock.Anything).Return(errorIn)
		boardService := &BoardService{boardStorage: boardStorage}
		err := boardService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
