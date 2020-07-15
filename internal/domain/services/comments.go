package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// CommentService is an interactor for work with comments
type CommentService struct {
	validator      v.Validator
	commentStorage CommentStorage
}

// NewCommentService is a comment service constructor
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

// Find will return all comments that meet the provided demand and an
// error in case it occurred while fetching records from the storage
func (c *CommentService) Find(demand CommentDemand) ([]*m.Comment, error) {
	return c.commentStorage.Find(demand)
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
