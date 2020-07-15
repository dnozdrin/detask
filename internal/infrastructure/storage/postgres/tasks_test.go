// +build unit

package postgres

import (
	"database/sql"
	"database/sql/driver"
	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTaskDAO_Save(t *testing.T) {
	t.Run("error_on_nil_task", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		tasksDAO := NewTaskDAO(db, logger)
		res, err := tasksDAO.Save(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("error_on_existing_ID", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		tasksDAO := NewTaskDAO(db, logger)
		task := &models.Task{Model: models.Model{ID: 1}}
		res, err := tasksDAO.Save(task)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, services.ErrRecordAlreadyExist, err)
	})
}

func TestTaskDAO_Update(t *testing.T) {
	t.Run("error_on_nil_task", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		tasksDAO := NewTaskDAO(db, logger)
		res, err := tasksDAO.Update(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Prepare", mock.Anything).Return(&sql.Stmt{}, errors.New("dummy"))
		tasksDAO := NewTaskDAO(db, logger)
		task := &models.Task{Model: models.Model{ID: 1}}
		res, err := tasksDAO.Update(task)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestTaskDAO_Delete(t *testing.T) {
	t.Run("exec_error", func(t *testing.T) {
		const ID uint = 0
		var result driver.RowsAffected = 0
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID}).Return(result, errors.New("dummy"))
		tasksDAO := NewTaskDAO(db, logger)
		err := tasksDAO.Delete(ID)

		assert.Error(t, err)
	})
}

func TestTaskDAO_MoveToColumn(t *testing.T) {
	t.Run("exec_error", func(t *testing.T) {
		const ID1, ID2 uint = 1, 2
		var result driver.RowsAffected = 0
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID2, ID1}).Return(result, errors.New("dummy"))
		tasksDAO := NewTaskDAO(db, logger)
		err := tasksDAO.MoveToColumn(ID1, ID2)

		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		const ID1, ID2 uint = 1, 2
		var result driver.RowsAffected = 0

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID2, ID1}).Return(result, nil)
		tasksDAO := NewTaskDAO(db, new(LoggerMock))
		err := tasksDAO.MoveToColumn(ID1, ID2)

		assert.Nil(t, err)
	})
}

func TestTaskDAO_WithTx(t *testing.T) {
	tx := &sql.Tx{}
	taskDAO := NewTaskDAO(new(QuerierMock), new(LoggerMock))
	txTaskDAO := taskDAO.WithTx(tx)

	assert.NotEqual(t, taskDAO, txTaskDAO)
	assert.Equal(t, txTaskDAO.(TaskDAO).db, tx)
}
