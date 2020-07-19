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

func TestCommentsDAO_Save(t *testing.T) {
	t.Run("error_on_nil_comment", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		commentsDAO := NewCommentsDAO(db, logger)
		res, err := commentsDAO.Save(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("error_on_existing_ID", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		commentsDAO := NewCommentsDAO(db, logger)
		comment := &models.Comment{Model: models.Model{ID: 1}}
		res, err := commentsDAO.Save(comment)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, services.ErrRecordAlreadyExist, err)
	})
}

func TestCommentsDAO_Update(t *testing.T) {
	t.Run("error_on_nil_comment", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		commentsDAO := NewCommentsDAO(db, logger)
		res, err := commentsDAO.Update(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Prepare", mock.Anything).Return(&sql.Stmt{}, errors.New("dummy"))
		commentsDAO := NewCommentsDAO(db, logger)
		comment := &models.Comment{Model: models.Model{ID: 1}}
		res, err := commentsDAO.Update(comment)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestCommentsDAO_Delete(t *testing.T) {
	t.Run("exec_error", func(t *testing.T) {
		const ID uint = 0
		var result driver.RowsAffected = 0
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID}).Return(result, errors.New("dummy"))
		commentsDAO := NewCommentsDAO(db, logger)
		err := commentsDAO.Delete(ID)

		assert.Error(t, err)
	})
}
