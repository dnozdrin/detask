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

func TestColumnDAO_Save(t *testing.T) {
	t.Run("error_on_nil_column", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		columnDAO := NewColumnDAO(db, logger)
		res, err := columnDAO.Save(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("error_on_existing_ID", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		columnDAO := NewColumnDAO(db, logger)
		column := &models.Column{Model: models.Model{ID: 1}}
		res, err := columnDAO.Save(column)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, services.ErrRecordAlreadyExist, err)
	})
}

func TestColumnDAO_Update(t *testing.T) {
	t.Run("error_on_nil_column", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Error", mock.Anything).Return()

		db := new(QuerierMock)
		columnDAO := NewColumnDAO(db, logger)
		res, err := columnDAO.Update(nil)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
	t.Run("prepare_error", func(t *testing.T) {
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Prepare", mock.Anything).Return(&sql.Stmt{}, errors.New("dummy"))
		columnDAO := NewColumnDAO(db, logger)
		column := &models.Column{Model: models.Model{ID: 1}}
		res, err := columnDAO.Update(column)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}

func TestColumnDAO_Delete(t *testing.T) {
	t.Run("exec_error", func(t *testing.T) {
		const ID uint = 0
		var result driver.RowsAffected = 0
		logger := new(LoggerMock)
		logger.On("Errorf", mock.Anything, mock.Anything).Return()

		db := new(QuerierMock)
		db.On("Exec", mock.Anything, []interface{}{ID}).Return(result, errors.New("dummy"))
		columnDAO := NewColumnDAO(db, logger)
		err := columnDAO.Delete(ID)

		assert.Error(t, err)
	})
}

func TestColumnDAO_WithTx(t *testing.T) {
	tx := &sql.Tx{}
	columnDAO := NewColumnDAO(new(QuerierMock), new(LoggerMock))
	txColumnDAO := columnDAO.WithTx(tx)

	assert.NotEqual(t, columnDAO, txColumnDAO)
	assert.Equal(t, txColumnDAO.(ColumnDAO).db, tx)
}
