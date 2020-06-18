package services

import (
	"errors"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskService(t *testing.T) {
	taskStorage := new(MockedTaskStorage)
	validation := new(MockedValidation)

	taskService := NewTaskService(taskStorage, validation)

	assert.Equal(t, taskStorage, taskService.taskStorage)
	assert.Equal(t, validation, taskService.validator)
}

func TestTaskService_Create(t *testing.T) {
	var taskIn = &m.Task{Name: "dummy"}
	var resultIn = v.NewResult(nil)
	t.Run("success", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Save", taskIn).Return(taskIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(resultIn)

		taskService := &TaskService{
			taskStorage: taskStorage,
			validator:   validation,
		}
		taskOut, resultOut := taskService.Create(taskIn)

		if assert.NotNil(t, taskOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})
	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(resultIn)

		taskService := &TaskService{validator: validation}
		taskOut, resultOut := taskService.Create(taskIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, taskOut)
	})
	t.Run("database_error", func(t *testing.T) {
		err := errors.New("simple error")

		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Save", taskIn).Return(&m.Task{}, err)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(v.NewResult(nil))

		taskService := NewTaskService(taskStorage, validation)
		taskOut, resultOut := taskService.Create(taskIn)

		assert.Equal(t, resultOut.Error, err)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_FindOneById(t *testing.T) {
	const dummyID = 1234
	taskIn := &m.Task{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindById", Anything).Return(taskIn, nil)
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, taskIn, taskOut)
	})

	t.Run("not_found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindById", Anything).Return(taskIn, errors.New(""))
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, taskIn, taskOut)
	})
}

func TestTaskService_FindAll(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		tasksIn := []*m.Task{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindAll").Return(tasksIn, nil)
		taskService := &TaskService{taskStorage: taskStorage}
		tasksOut, err := taskService.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, tasksIn, tasksOut)
	})

	t.Run("not_found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindAll", Anything).Return([]*m.Task{}, errors.New(""))
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.FindAll()
		assert.Error(t, err)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_Update(t *testing.T) {
	var taskIn = &m.Task{Name: "dummy"}
	var resultIn = v.NewResult(nil)

	t.Run("success", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Update", taskIn).Return(taskIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(resultIn)

		taskService := &TaskService{
			taskStorage: taskStorage,
			validator:   validation,
		}
		taskOut, resultOut := taskService.Update(taskIn)

		if assert.NotNil(t, taskOut) {
			assert.Equal(t, resultIn, resultOut)
		}
	})

	t.Run("validation_error", func(t *testing.T) {
		resultIn = v.NewResult(v.ErrValidationFailed)
		resultIn.Errors = append(resultIn.Errors, v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(resultIn)

		taskService := &TaskService{validator: validation}
		taskOut, resultOut := taskService.Update(taskIn)

		assert.Equal(t, resultIn, resultOut)
		assert.Empty(t, taskOut)
	})

	t.Run("database_error", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Update", taskIn).Return(&m.Task{}, errors.New("simple error"))

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(v.NewResult(nil))

		taskService := NewTaskService(taskStorage, validation)
		taskOut, resultOut := taskService.Update(taskIn)

		assert.Error(t, resultOut.Error)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Delete", Anything).Return(nil)
		taskService := &TaskService{taskStorage: taskStorage}
		err := taskService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Delete", Anything).Return(errorIn)
		taskService := &TaskService{taskStorage: taskStorage}
		err := taskService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
