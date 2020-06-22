package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// ColumnStorage represents an interface for interaction with columns DAO
type ColumnStorage interface {
	// Save will persist the provided column
	Save(column *m.Column) (*m.Column, error)
	// FindOneById should return a column with the provided ID
	FindOneById(ID uint) (*m.Column, error)
	// Find should return a slice of columns pointers sorted by position, that meet the
	// provided demand
	Find() ([]*m.Column, error)
	// Update should update all column fields by the provided data
	Update(column *m.Column) (*m.Column, error)
	// Delete should set current deletion time to a column with the provided ID
	// and to all dependant records
	Delete(ID uint) error
}

// ColumnService is an interactor for work with columns
type ColumnService struct {
	validator     v.Validator
	columnStorage ColumnStorage
}

// NewColumnService is a column service constructor
func NewColumnService(validator v.Validator, columnStorage ColumnStorage) *ColumnService {
	return &ColumnService{
		columnStorage: columnStorage,
		validator:     validator,
	}
}

// Create will create a new column with the provided payload. Returns the
// operation result with possible validation or saving errors
func (c *ColumnService) Create(column *m.Column) (*m.Column, error) {
	if err := c.validator.Validate(*column); err != nil {
		return nil, err
	}

	return c.columnStorage.Save(column)
}

// Find will return all not deleted columns and an error in case
// it occurred while fetching records from the storage
func (c *ColumnService) Find() ([]*m.Column, error) {
	return c.columnStorage.Find()
}

// FindOneById will return a pointer to the column requested by id and
// an error in case it occurred while fetching the record from the storage
func (c *ColumnService) FindOneById(ID uint) (*m.Column, error) {
	return c.columnStorage.FindOneById(ID)
}

// Update will update the column record. Returns the operation result
// with possible validation or saving errors
func (c *ColumnService) Update(column *m.Column) (*m.Column, error) {
	if err := c.validator.Validate(*column); err != nil {
		return nil, err
	}

	return c.columnStorage.Update(column)
}

// Delete will mark a record with the given ID as deleted as well as all
// the dependant records
func (c *ColumnService) Delete(ID uint) error {
	return c.columnStorage.Delete(ID)
}
