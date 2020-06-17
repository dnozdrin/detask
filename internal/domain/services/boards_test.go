package services

import (
	"errors"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardService(t *testing.T) {
	boardStorage := new(MockedBoardStorage)
	validation := new(MockedValidation)

	boardService := NewBoardService(boardStorage, validation)

	assert.Equal(t, boardStorage, boardService.boardStorage)
	assert.Equal(t, validation, boardService.validator)
}

func TestBoardService_Create(t *testing.T) {
	var boardIn = &m.Board{Name: "dummy"}
	var resultIn = v.NewResult(nil)
	t.Run("success", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(boardIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(resultIn)

		boardService := &BoardService{
			boardStorage: boardStorage,
			validator:    validation,
		}
		boardOut, resultOut := boardService.Create(boardIn)

		if assert.NotNil(t, boardOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})
	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(resultIn)

		boardService := &BoardService{validator: validation}
		boardOut, resultOut := boardService.Create(boardIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, boardOut)
	})
	t.Run("database_error", func(t *testing.T) {
		err := errors.New("simple error")

		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Save", boardIn).Return(&m.Board{}, err)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(v.NewResult(nil))

		boardService := NewBoardService(boardStorage, validation)
		boardOut, resultOut := boardService.Create(boardIn)

		assert.Equal(t, resultOut.Error, err)
		assert.Empty(t, boardOut)
	})
}

func TestBoardService_FindOneById(t *testing.T) {
	const dummyID = 1234
	boardIn := &m.Board{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindById", Anything).Return(boardIn, nil)
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, boardIn, boardOut)
	})

	t.Run("not_found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindById", Anything).Return(boardIn, errors.New(""))
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, boardIn, boardOut)
	})
}

func TestBoardService_FindAll(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		boardsIn := []*m.Board{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindAll").Return(boardsIn, nil)
		boardService := &BoardService{boardStorage: boardStorage}
		boardsOut, err := boardService.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, boardsIn, boardsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("FindAll", Anything).Return([]*m.Board{}, errors.New(""))
		boardService := &BoardService{boardStorage: boardStorage}
		boardOut, err := boardService.FindAll()
		assert.Error(t, err)
		assert.Empty(t, boardOut)
	})
}

func TestBoardService_Update(t *testing.T) {
	var boardIn = &m.Board{Name: "dummy"}
	var resultIn = v.NewResult(nil)

	t.Run("success", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Update", boardIn).Return(boardIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(resultIn)

		boardService := &BoardService{
			boardStorage: boardStorage,
			validator:    validation,
		}
		boardOut, resultOut := boardService.Update(boardIn)

		if assert.NotNil(t, boardOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})

	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(resultIn)

		boardService := &BoardService{validator: validation}
		boardOut, resultOut := boardService.Update(boardIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, boardOut)
	})

	t.Run("database_error", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Update", boardIn).Return(&m.Board{}, errors.New("simple error"))

		validation := new(MockedValidation)
		validation.On("Validate", *boardIn).Return(v.NewResult(nil))

		boardService := NewBoardService(boardStorage, validation)
		boardOut, resultOut := boardService.Update(boardIn)

		assert.Error(t, resultOut.Error)
		assert.Empty(t, boardOut)
	})
}

func TestBoardService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Delete", Anything).Return(nil)
		boardService := &BoardService{boardStorage: boardStorage}
		err := boardService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		boardStorage := new(MockedBoardStorage)
		boardStorage.On("Delete", Anything).Return(errorIn)
		boardService := &BoardService{boardStorage: boardStorage}
		err := boardService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
