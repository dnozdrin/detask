package services

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewCommentService(t *testing.T) {
	commentStorage := new(MockedCommentStorage)
	validation := new(MockedValidation)
	commentService := NewCommentService(validation, commentStorage)

	assert.Equal(t, commentStorage, commentService.commentStorage)
	assert.Equal(t, validation, commentService.validator)
}

func TestCommentService_Create(t *testing.T) {
	var commentIn = &m.Comment{Text: "dummy"}
	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Save", commentIn).Return(commentIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}
		commentOut, err := commentService.Create(commentIn)

		assert.NotNil(t, commentOut)
		assert.Nil(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{validator: validation}
		commentOut, err := commentService.Create(commentIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, commentOut)
	})
	t.Run("database_error", func(t *testing.T) {
		var validationErr *v.Errors
		dbErr := errors.New("simple error")
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Save", commentIn).Return(&m.Comment{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}

		commentOut, err := commentService.Create(commentIn)

		assert.Empty(t, commentOut)
		assert.Equal(t, err, dbErr)
	})
}

func TestCommentService_FindOneById(t *testing.T) {
	const dummyID = 1234
	commentIn := &m.Comment{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindOneById", mock.Anything).Return(commentIn, nil)
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, commentIn, commentOut)
	})

	t.Run("not_found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindOneById", mock.Anything).Return(commentIn, errors.New(""))
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, commentIn, commentOut)
	})
}

func TestCommentService_Find(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		commentsIn := []*m.Comment{
			{Text: "Test1"},
			{Text: "Test2"},
		}
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Find", mock.Anything).Return(commentsIn, nil)
		commentService := &CommentService{commentStorage: commentStorage}
		commentsOut, err := commentService.Find(make(CommentDemand))
		assert.Nil(t, err)
		assert.Equal(t, commentsIn, commentsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Find", mock.Anything).Return([]*m.Comment{}, errors.New(""))
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.Find(make(CommentDemand))
		assert.Error(t, err)
		assert.Empty(t, commentOut)
	})
}

func TestCommentService_Update(t *testing.T) {
	var commentIn = &m.Comment{Text: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Update", commentIn).Return(commentIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}
		commentOut, err := commentService.Update(commentIn)

		assert.NotNil(t, commentOut)
		assert.Nil(t, err)
	})

	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{validator: validation}
		commentOut, err := commentService.Update(commentIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, commentOut)
	})

	t.Run("database_error", func(t *testing.T) {
		dbErr := errors.New("simple error")
		var validationErr *v.Errors
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Update", commentIn).Return(&m.Comment{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(validationErr)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}
		commentOut, err := commentService.Update(commentIn)

		assert.Empty(t, commentOut)
		assert.Equal(t, err, dbErr)
	})
}

func TestCommentService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Delete", mock.Anything).Return(nil)
		commentService := &CommentService{commentStorage: commentStorage}
		err := commentService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Delete", mock.Anything).Return(errorIn)
		commentService := &CommentService{commentStorage: commentStorage}
		err := commentService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
