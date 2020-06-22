package services

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardService(t *testing.T) {
	boardStorage := new(MockedBoardStorage)
	validation := new(MockedValidation)
	boardService := NewBoardService(validation, boardStorage)

	assert.Equal(t, validation, boardService.validator)
	assert.Equal(t, boardStorage, boardService.boardStorage)
}

func TestBoardService_Create(t *testing.T) {
	var boardIn = &m.Board{Name: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("SaveWithDefaultColumn", boardIn).Return(boardIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{
			validator:    validation,
			boardStorage: boardStorage,
		}
		boardOut, err := boardService.Create(boardIn)

		assert.NotNil(t, boardOut)
		assert.Nil(t, err)
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
	t.Run("database_error", func(t *testing.T) {
		var validationErr *v.Errors
		dbErr := errors.New("simple error")

		boardStorage := new(MockedBoardStorage)
		boardStorage.On("SaveWithDefaultColumn", boardIn).Return(&m.Board{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(validationErr)

		boardService := &BoardService{
			validator:    validation,
			boardStorage: boardStorage,
		}
		boardOut, err := boardService.Create(boardIn)

		assert.Empty(t, boardOut)
		assert.Equal(t, err, dbErr)
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
