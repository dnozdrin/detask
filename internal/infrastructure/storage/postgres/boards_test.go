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

func TestBoardDAO_Save(t *testing.T) {
	t.Run("error_on_nil_board", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		boardDAO := NewBoardDAO(db, logger)
		res, err := boardDAO.Save(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("error_on_existing_ID", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		boardDAO := NewBoardDAO(db, logger)
		board := &models.Board{Model: models.Model{ID: 1}}
		res, err := boardDAO.Save(board)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, services.ErrRecordAlreadyExist, err)
	})
	t.Run("stmt_prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Prepare", mock.Anything).Return(&sql.Stmt{}, errors.New("dummy"))
		boardDAO := NewBoardDAO(db, logger)
		board := &models.Board{Model: models.Model{}}
		res, err := boardDAO.Save(board)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestBoardDAO_Find(t *testing.T) {
	t.Run("query_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Query", mock.Anything, mock.Anything).Return(&sql.Rows{}, errors.New("dummy"))
		boardDAO := NewBoardDAO(db, logger)
		res, err := boardDAO.Find()

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestBoardDAO_Update(t *testing.T) {
	t.Run("error_on_nil_board", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		boardDAO := NewBoardDAO(db, logger)
		res, err := boardDAO.Update(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Prepare", mock.Anything).Return(&sql.Stmt{}, errors.New("dummy"))
		boardDAO := NewBoardDAO(db, logger)
		board := &models.Board{Model: models.Model{ID: 1}}
		res, err := boardDAO.Update(board)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestBoardDAO_Delete(t *testing.T) {
	t.Run("exec_error", func(t *testing.T) {
		const ID uint = 0
		var result driver.RowsAffected = 0
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID}).Return(result, errors.New("dummy"))
		boardDAO := NewBoardDAO(db, logger)
		err := boardDAO.Delete(ID)

		assert.Error(t, err)
	})
}

func TestBoardDAO_WithTx(t *testing.T) {
	tx := &sql.Tx{}
	boardDAO := NewBoardDAO(new(QuerierMock), new(LoggerMock))
	txBoardDAO := boardDAO.WithTx(tx)

	assert.NotEqual(t, boardDAO, txBoardDAO)
	assert.Equal(t, txBoardDAO.(BoardDAO).db, tx)
}
