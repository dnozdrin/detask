package services

import (
	"errors"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewCommentService(t *testing.T) {
	commentStorage := new(MockedCommentStorage)
	validation := new(MockedValidation)

	commentService := NewCommentService(commentStorage, validation)

	assert.Equal(t, commentStorage, commentService.commentStorage)
	assert.Equal(t, validation, commentService.validator)
}

func TestCommentService_Create(t *testing.T) {
	var commentIn = &m.Comment{Text: "dummy"}
	var resultIn = v.NewResult(nil)
	t.Run("success", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Save", commentIn).Return(commentIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(resultIn)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}
		commentOut, resultOut := commentService.Create(commentIn)

		if assert.NotNil(t, commentOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})
	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(resultIn)

		commentService := &CommentService{validator: validation}
		commentOut, resultOut := commentService.Create(commentIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, commentOut)
	})
	t.Run("database_error", func(t *testing.T) {
		err := errors.New("simple error")

		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Save", commentIn).Return(&m.Comment{}, err)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(v.NewResult(nil))

		commentService := NewCommentService(commentStorage, validation)
		commentOut, resultOut := commentService.Create(commentIn)

		assert.Equal(t, resultOut.Error, err)
		assert.Empty(t, commentOut)
	})
}

func TestCommentService_FindOneById(t *testing.T) {
	const dummyID = 1234
	commentIn := &m.Comment{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindById", Anything).Return(commentIn, nil)
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, commentIn, commentOut)
	})

	t.Run("not_found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindById", Anything).Return(commentIn, errors.New(""))
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, commentIn, commentOut)
	})
}

func TestCommentService_FindAll(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		commentsIn := []*m.Comment{
			{Text: "Test1"},
			{Text: "Test2"},
		}
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindAll").Return(commentsIn, nil)
		commentService := &CommentService{commentStorage: commentStorage}
		commentsOut, err := commentService.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, commentsIn, commentsOut)
	})

	t.Run("not_found", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("FindAll", Anything).Return([]*m.Comment{}, errors.New(""))
		commentService := &CommentService{commentStorage: commentStorage}
		commentOut, err := commentService.FindAll()
		assert.Error(t, err)
		assert.Empty(t, commentOut)
	})
}

func TestCommentService_Update(t *testing.T) {
	var commentIn = &m.Comment{Text: "dummy"}
	var resultIn = v.NewResult(nil)

	t.Run("success", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Update", commentIn).Return(commentIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(resultIn)

		commentService := &CommentService{
			commentStorage: commentStorage,
			validator:      validation,
		}
		commentOut, resultOut := commentService.Update(commentIn)

		if assert.NotNil(t, commentOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})

	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(resultIn)

		commentService := &CommentService{validator: validation}
		commentOut, resultOut := commentService.Update(commentIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, commentOut)
	})

	t.Run("database_error", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Update", commentIn).Return(&m.Comment{}, errors.New("simple error"))

		validation := new(MockedValidation)
		validation.On("Validate", *commentIn).Return(v.NewResult(nil))

		commentService := NewCommentService(commentStorage, validation)
		commentOut, resultOut := commentService.Update(commentIn)

		assert.Error(t, resultOut.Error)
		assert.Empty(t, commentOut)
	})
}

func TestCommentService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Delete", Anything).Return(nil)
		commentService := &CommentService{commentStorage: commentStorage}
		err := commentService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		commentStorage := new(MockedCommentStorage)
		commentStorage.On("Delete", Anything).Return(errorIn)
		commentService := &CommentService{commentStorage: commentStorage}
		err := commentService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
