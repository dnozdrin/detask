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

		db := new(DBMock)
		tasksDAO := NewTaskDAO(db, logger)
		res, err := tasksDAO.Save(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("error_on_existing_ID", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		db := new(DBMock)
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

		db := new(DBMock)
		tasksDAO := NewTaskDAO(db, logger)
		res, err := tasksDAO.Update(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(DBMock)
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

		db := new(DBMock)
		db.On("Exec", mock.Anything, []interface{}{ID}).Return(result, errors.New("dummy"))
		tasksDAO := NewTaskDAO(db, logger)
		err := tasksDAO.Delete(ID)

		assert.Error(t, err)
	})
}
