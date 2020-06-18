package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

type CommentStorage interface {
	Save(comment *m.Comment) (*m.Comment, error)
	FindById(id uint) (*m.Comment, error)
	FindAll() ([]*m.Comment, error)
	Update(comment *m.Comment) (*m.Comment, error)
	Delete(id uint) error
}

type CommentService struct {
	commentStorage CommentStorage
	validator      v.Validator
}

// CommentService is a comment service constructor
func NewCommentService(commentStorage CommentStorage, validator v.Validator) *CommentService {
	return &CommentService{
		commentStorage: commentStorage,
		validator:      validator,
	}
}

// Create will create a new comment  with the provided payload. Returns the
// operation result with possible validation or saving errors
func (c *CommentService) Create(comment *m.Comment) (*m.Comment, v.Result) {
	var result v.Result
	if result = c.validator.Validate(*comment); !result.IsValid() {
		return nil, result
	}

	comment, result.Error = c.commentStorage.Save(comment)

	return comment, result
}

// FindAll will return all not deleted comments and an error in case
// it occurred while fetching records from the storage
func (c *CommentService) FindAll() ([]*m.Comment, error) {
	return c.commentStorage.FindAll()
}

// FindOneById will return a pointer to the comment requested by id and
// an error in case it occurred while fetching the record from the storage
func (c *CommentService) FindOneById(ID uint) (*m.Comment, error) {
	return c.commentStorage.FindById(ID)
}

// Update will update the comment record. Returns the operation result
// with possible validation or saving errors
func (c *CommentService) Update(comment *m.Comment) (*m.Comment, v.Result) {
	var result v.Result
	if result = c.validator.Validate(*comment); !result.IsValid() {
		return nil, result
	}

	comment, result.Error = c.commentStorage.Update(comment)

	return comment, result
}

// Delete will delete a record with the given ID
func (c *CommentService) Delete(ID uint) error {
	return c.commentStorage.Delete(ID)
}
