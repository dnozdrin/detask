package services

import (
	"errors"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewColumnService(t *testing.T) {
	columnStorage := new(MockedColumnStorage)
	validation := new(MockedValidation)

	columnService := NewColumnService(columnStorage, validation)

	assert.Equal(t, columnStorage, columnService.columnStorage)
	assert.Equal(t, validation, columnService.validator)
}

func TestColumnService_Create(t *testing.T) {
	var columnIn = &m.Column{Name: "dummy"}
	var resultIn = v.NewResult(nil)
	t.Run("success", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", columnIn).Return(columnIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(resultIn)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			validator:     validation,
		}
		columnOut, resultOut := columnService.Create(columnIn)

		if assert.NotNil(t, columnOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})
	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(resultIn)

		columnService := &ColumnService{validator: validation}
		columnOut, resultOut := columnService.Create(columnIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, columnOut)
	})
	t.Run("database_error", func(t *testing.T) {
		err := errors.New("simple error")

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", columnIn).Return(&m.Column{}, err)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(v.NewResult(nil))

		columnService := NewColumnService(columnStorage, validation)
		columnOut, resultOut := columnService.Create(columnIn)

		assert.Equal(t, resultOut.Error, err)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_FindOneById(t *testing.T) {
	const dummyID = 1234
	columnIn := &m.Column{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindById", Anything).Return(columnIn, nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, columnIn, columnOut)
	})

	t.Run("not_found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindById", Anything).Return(columnIn, errors.New(""))
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, columnIn, columnOut)
	})
}

func TestColumnService_FindAll(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		columnsIn := []*m.Column{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindAll").Return(columnsIn, nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		columnsOut, err := columnService.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, columnsIn, columnsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindAll", Anything).Return([]*m.Column{}, errors.New(""))
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.FindAll()
		assert.Error(t, err)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_Update(t *testing.T) {
	var columnIn = &m.Column{Name: "dummy"}
	var resultIn = v.NewResult(nil)

	t.Run("success", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Update", columnIn).Return(columnIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(resultIn)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			validator:     validation,
		}
		columnOut, resultOut := columnService.Update(columnIn)

		if assert.NotNil(t, columnOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})

	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(resultIn)

		columnService := &ColumnService{validator: validation}
		columnOut, resultOut := columnService.Update(columnIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, columnOut)
	})

	t.Run("database_error", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Update", columnIn).Return(&m.Column{}, errors.New("simple error"))

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(v.NewResult(nil))

		columnService := NewColumnService(columnStorage, validation)
		columnOut, resultOut := columnService.Update(columnIn)

		assert.Error(t, resultOut.Error)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", Anything).Return(nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		err := columnService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", Anything).Return(errorIn)
		columnService := &ColumnService{columnStorage: columnStorage}
		err := columnService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
