package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// DefaultColPos is the default position for columns that are created internally
const DefaultColPos = 1000

// ColumnService is an interactor for work with columns
type ColumnService struct {
	validator     v.Validator
	columnStorage ColumnStorage
	taskStorage   TaskStorage
	txBeginner    TxBeginner
}

// NewColumnService is a column service constructor
func NewColumnService(
	validator v.Validator,
	columnStorage ColumnStorage,
	taskStorage TaskStorage,
	txBeginner TxBeginner,
) ColumnService {
	return ColumnService{
		columnStorage: columnStorage,
		taskStorage:   taskStorage,
		validator:     validator,
		txBeginner:    txBeginner,
	}
}

// Create will create a new column with the provided payload. Returns the
// operation result with possible validation or saving errors
func (c ColumnService) Create(column *m.Column) (*m.Column, error) {
	if err := c.validator.Validate(*column); err != nil {
		return nil, err
	}

	return c.columnStorage.Save(column)
}

// Find will return all not deleted columns and an error in case
// it occurred while fetching records from the storage
func (c ColumnService) Find(demand ColumnDemand) ([]*m.Column, error) {
	return c.columnStorage.Find(demand)
}

// FindOneById will return a pointer to the column requested by id and
// an error in case it occurred while fetching the record from the storage
func (c ColumnService) FindOneById(ID uint) (*m.Column, error) {
	return c.columnStorage.FindOneById(ID)
}

// Update will update the column record. Returns the operation result
// with possible validation or saving errors
func (c ColumnService) Update(column *m.Column) (*m.Column, error) {
	if err := c.validator.Validate(*column); err != nil {
		return nil, err
	}

	return c.columnStorage.Update(column)
}

// Delete will the column with the provided ID. The last column cannot be deleted.
// When a column is deleted, its tasks are moved to the column to the left of the
// current or to the right of the current if the curring is the leftmost
func (c ColumnService) Delete(ID uint) error {
	tx, err := c.txBeginner.Begin()
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()
	columnStorage := c.columnStorage.WithTx(tx)
	column, err := columnStorage.FindOneById(ID)
	if err != nil {
		return ErrRecordNotFound
	}

	columnsNum, err := columnStorage.CountColumnsByBoard(column.BoardID)
	if err != nil {
		return err
	}
	if columnsNum == 1 {
		return ErrLastColumn
	}

	targetColumn, err := columnStorage.FindColumnToTheLeft(ID)
	if err != nil {
		targetColumn, err = columnStorage.FindColumnToTheRight(ID)
		if err != nil {
			return ErrTargetColumn
		}
	}
	taskStorage := c.taskStorage.WithTx(tx)
	if err = taskStorage.MoveToColumn(ID, targetColumn); err != nil {
		return err
	}

	if err = columnStorage.Delete(ID); err != nil {
		return err
	}

	return tx.Commit()
}
