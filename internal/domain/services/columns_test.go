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

func TestNewColumnService(t *testing.T) {
	columnStorage := new(MockedColumnStorage)
	validation := new(MockedValidation)
	txBeginner := new(MockedTxBeginner)
	taskStorage := new(MockedTaskStorage)
	columnService := NewColumnService(validation, columnStorage, taskStorage, txBeginner)

	assert.Equal(t, columnStorage, columnService.columnStorage)
	assert.Equal(t, taskStorage, columnService.taskStorage)
	assert.Equal(t, txBeginner, columnService.txBeginner)
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
	const (
		currColID      uint = 123
		boardId        uint = 12
		leftColID      uint = 24
		rightColID     uint = 70
		columnsOnBoard int  = 5
	)

	t.Run("successful_delete", func(t *testing.T) {
		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("WithTx", tx).Return(taskStorage)
		taskStorage.On("MoveToColumn", currColID, leftColID).Return(nil)

		currColumn := &m.Column{Name: "Test", Model: m.Model{ID: currColID}, BoardID: boardId}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(currColumn, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(columnsOnBoard, nil)
		columnStorage.On("FindColumnToTheLeft", currColID).Return(leftColID, nil)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			taskStorage:   taskStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Nil(t, err)
	})

	t.Run("last_column", func(t *testing.T) {
		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectRollback()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("WithTx", tx).Return(taskStorage)
		taskStorage.On("MoveToColumn", currColID, leftColID).Return(nil)

		currColumn := &m.Column{Name: "Test", Model: m.Model{ID: currColID}, BoardID: boardId}
		columnStorage := new(MockedColumnStorage)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(currColumn, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(1, nil)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, ErrLastColumn, err)
	})

	t.Run("tx_begin_error", func(t *testing.T) {
		txErr := errors.New("tx error")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, txErr)

		columnService := &ColumnService{
			txBeginner: txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, txErr, err)
	})

	t.Run("current_record_search_error", func(t *testing.T) {
		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectRollback()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{}, errors.New("not found"))

		columnService := &ColumnService{
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, ErrRecordNotFound, err)
	})

	t.Run("columns_count_error", func(t *testing.T) {
		countErr := errors.New("count error")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectRollback()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{BoardID: boardId}, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(0, countErr)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, countErr, err)
	})

	t.Run("left_column_search_error", func(t *testing.T) {
		searchErr := errors.New("left not found")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("WithTx", tx).Return(taskStorage)
		taskStorage.On("MoveToColumn", currColID, rightColID).Return(nil)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{BoardID: boardId}, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(columnsOnBoard, nil)
		columnStorage.On("FindColumnToTheLeft", currColID).Return(uint(0), searchErr)
		columnStorage.On("FindColumnToTheRight", currColID).Return(rightColID, nil)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			taskStorage:   taskStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Nil(t, err)
	})

	t.Run("tasks_target_search_error", func(t *testing.T) {
		searchErr := errors.New("left not found")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{BoardID: boardId}, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(columnsOnBoard, nil)
		columnStorage.On("FindColumnToTheLeft", currColID).Return(uint(0), searchErr)
		columnStorage.On("FindColumnToTheRight", currColID).Return(uint(0), searchErr)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, ErrTargetColumn, err)
	})

	t.Run("tasks_move_error", func(t *testing.T) {
		moveErr := errors.New("error on tasks move")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("WithTx", tx).Return(taskStorage)
		taskStorage.On("MoveToColumn", currColID, leftColID).Return(moveErr)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(nil)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{BoardID: boardId}, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(columnsOnBoard, nil)
		columnStorage.On("FindColumnToTheLeft", currColID).Return(leftColID, nil)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			taskStorage:   taskStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, moveErr, err)
	})

	t.Run("deletion_error", func(t *testing.T) {
		dbErr := errors.New("deletion error")

		db, dbmock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		dbmock.ExpectBegin()
		dbmock.ExpectCommit()
		tx, _ := db.Begin()

		txBeginner := new(MockedTxBeginner)
		txBeginner.On("Begin").Return(tx, nil)

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("WithTx", tx).Return(taskStorage)
		taskStorage.On("MoveToColumn", currColID, leftColID).Return(nil)

		columnStorage := new(MockedColumnStorage)
		columnStorage.On("Delete", currColID).Return(dbErr)
		columnStorage.On("WithTx", tx).Return(columnStorage)
		columnStorage.On("FindOneById", currColID).Return(&m.Column{BoardID: boardId}, nil)
		columnStorage.On("CountColumnsByBoard", boardId).Return(columnsOnBoard, nil)
		columnStorage.On("FindColumnToTheLeft", currColID).Return(leftColID, nil)

		columnService := &ColumnService{
			columnStorage: columnStorage,
			taskStorage:   taskStorage,
			txBeginner:    txBeginner,
		}
		err = columnService.Delete(currColID)
		assert.Equal(t, dbErr, err)
	})
}
