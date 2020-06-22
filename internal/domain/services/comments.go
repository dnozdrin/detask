package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// CommentStorage represents an interface for interaction with comments DAO
type CommentStorage interface {
	// Save will persist the provided comment
	Save(comment *m.Comment) (*m.Comment, error)
	// FindOneById should return a comment with the provided ID
	FindOneById(id uint) (*m.Comment, error)
	// Find should return a slice of comments pointers sorted by creation date
	// (from newest to oldest), that meet the provided demand
	Find() ([]*m.Comment, error)
	// Update should update the comment text
	Update(comment *m.Comment) (*m.Comment, error)
	// Delete should set current deletion time to a comment with the provided ID
	// and to all dependant records
	Delete(id uint) error
}

// ColumnService is an interactor for work with comments
type CommentService struct {
	validator      v.Validator
	commentStorage CommentStorage
}

// CommentService is a comment service constructor
func NewCommentService(validator v.Validator, commentStorage CommentStorage) *CommentService {
	return &CommentService{
		commentStorage: commentStorage,
		validator:      validator,
	}
}

// Create will create a new comment  with the provided payload. Returns the
// operation result with possible validation or saving errors
func (c *CommentService) Create(comment *m.Comment) (*m.Comment, error) {
	if err := c.validator.Validate(*comment); err != nil {
		return nil, err
	}

	return c.commentStorage.Save(comment)
}

// Find will return all not deleted comments and an error in case
// it occurred while fetching records from the storage
func (c *CommentService) Find() ([]*m.Comment, error) {
	return c.commentStorage.Find()
}

// FindOneById will return a pointer to the comment requested by id and
// an error in case it occurred while fetching the record from the storage
func (c *CommentService) FindOneById(ID uint) (*m.Comment, error) {
	return c.commentStorage.FindOneById(ID)
}

// Update will update the comment record. Returns the operation result
// with possible validation or saving errors
func (c *CommentService) Update(comment *m.Comment) (*m.Comment, error) {
	if err := c.validator.Validate(*comment); err != nil {
		return nil, err
	}

	return c.commentStorage.Update(comment)
}

// Delete will delete a record with the given ID
func (c *CommentService) Delete(ID uint) error {
	return c.commentStorage.Delete(ID)
}
