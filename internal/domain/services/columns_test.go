// +build unit

package services

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewColumnService(t *testing.T) {
	columnStorage := new(MockedColumnStorage)
	validation := new(MockedValidation)

	columnService := NewColumnService(validation, columnStorage)

	assert.Equal(t, columnStorage, columnService.columnStorage)
	assert.Equal(t, validation, columnService.validator)
}

func TestColumnService_Create(t *testing.T) {
	var columnIn = &m.Column{Name: "dummy"}
	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", columnIn).Return(columnIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{
			validator:     validation,
			columnStorage: columnStorage,
		}
		columnOut, err := columnService.Create(columnIn)

		assert.NotNil(t, columnOut)
		assert.Nil(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{validator: validation}
		columnOut, err := columnService.Create(columnIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, columnOut)
	})
	t.Run("database_error", func(t *testing.T) {
		var validationErr *v.Errors
		err := errors.New("simple error")

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Save", columnIn).Return(&m.Column{}, err)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{
			validator:     validation,
			columnStorage: columnStorage,
		}
		columnOut, resultOut := columnService.Create(columnIn)

		assert.Equal(t, resultOut, err)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_FindOneById(t *testing.T) {
	const dummyID = 1234
	columnIn := &m.Column{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindOneById", mock.Anything).Return(columnIn, nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, columnIn, columnOut)
	})

	t.Run("not_found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("FindOneById", mock.Anything).Return(columnIn, errors.New(""))
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, columnIn, columnOut)
	})
}

func TestColumnService_Find(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		columnsIn := []*m.Column{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Find", mock.Anything).Return(columnsIn, nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		columnsOut, err := columnService.Find(make(ColumnDemand))
		assert.Nil(t, err)
		assert.Equal(t, columnsIn, columnsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Find", mock.Anything).Return([]*m.Column{}, errors.New(""))
		columnService := &ColumnService{columnStorage: columnStorage}
		columnOut, err := columnService.Find(make(ColumnDemand))
		assert.Error(t, err)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_Update(t *testing.T) {
	var columnIn = &m.Column{Name: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Update", columnIn).Return(columnIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			validator:     validation,
		}
		columnOut, resultOut := columnService.Update(columnIn)

		assert.NotNil(t, columnOut)
		assert.Nil(t, resultOut)
	})

	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{validator: validation}
		columnOut, err := columnService.Update(columnIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, columnOut)
	})

	t.Run("database_error", func(t *testing.T) {
		var validationErr *v.Errors
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Update", columnIn).Return(&m.Column{}, errors.New("simple error"))

		validation := new(MockedValidation)
		validation.On("Validate", *columnIn).Return(validationErr)

		columnService := &ColumnService{
			validator:     validation,
			columnStorage: columnStorage,
		}
		columnOut, resultOut := columnService.Update(columnIn)

		assert.Error(t, resultOut)
		assert.Empty(t, columnOut)
	})
}

func TestColumnService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", mock.Anything).Return(nil)
		columnService := &ColumnService{columnStorage: columnStorage}
		err := columnService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", mock.Anything).Return(errorIn)
		columnService := &ColumnService{columnStorage: columnStorage}
		err := columnService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
