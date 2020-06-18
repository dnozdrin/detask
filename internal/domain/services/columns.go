package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

type ColumnStorage interface {
	Save(column *m.Column) (*m.Column, error)
	FindById(id uint) (*m.Column, error)
	FindAll() ([]*m.Column, error)
	Update(column *m.Column) (*m.Column, error)
	Delete(id uint) error
}

type ColumnService struct {
	columnStorage ColumnStorage
	validator     v.Validator
}

// NewColumnService is a column service constructor
func NewColumnService(columnStorage ColumnStorage, validator v.Validator) *ColumnService {
	return &ColumnService{
		columnStorage: columnStorage,
		validator:     validator,
	}
}

// Create will create a new column with the provided payload. Returns the
// operation result with possible validation or saving errors
func (c *ColumnService) Create(column *m.Column) (*m.Column, v.Result) {
	var result v.Result
	if result = c.validator.Validate(*column); !result.IsValid() {
		return nil, result
	}

	column, result.Error = c.columnStorage.Save(column)

	return column, result
}

// FindAll will return all not deleted columns and an error in case
// it occurred while fetching records from the storage
func (c *ColumnService) FindAll() ([]*m.Column, error) {
	return c.columnStorage.FindAll()
}

// FindOneById will return a pointer to the column requested by id and
// an error in case it occurred while fetching the record from the storage
func (c *ColumnService) FindOneById(ID uint) (*m.Column, error) {
	return c.columnStorage.FindById(ID)
}

// Update will update the column record. Returns the operation result
// with possible validation or saving errors
func (c *ColumnService) Update(column *m.Column) (*m.Column, v.Result) {
	var result v.Result
	if result = c.validator.Validate(*column); !result.IsValid() {
		return nil, result
	}

	column, result.Error = c.columnStorage.Update(column)

	return column, result
}

// Delete will delete a record with the given ID
func (c *ColumnService) Delete(ID uint) error {
	return c.columnStorage.Delete(ID)
}
