// +build unit

package services

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"

	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskService(t *testing.T) {
	taskStorage := new(MockedTaskStorage)
	validation := new(MockedValidation)
	taskService := NewTaskService(validation, taskStorage)

	assert.Equal(t, validation, taskService.validator)
	assert.Equal(t, taskStorage, taskService.taskStorage)
}

func TestTaskService_Create(t *testing.T) {
	var taskIn = &m.Task{Name: "dummy"}
	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Save", taskIn).Return(taskIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{
			taskStorage: taskStorage,
			validator:   validation,
		}
		taskOut, err := taskService.Create(taskIn)

		assert.NotNil(t, taskOut)
		assert.Nil(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{validator: validation}
		taskOut, err := taskService.Create(taskIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, taskOut)
	})
	t.Run("database_error", func(t *testing.T) {
		dbErr := errors.New("simple error")
		var validationErr *v.Errors
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Save", taskIn).Return(&m.Task{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{
			taskStorage: taskStorage,
			validator:   validation,
		}
		taskOut, err := taskService.Create(taskIn)

		assert.Equal(t, dbErr, err)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_FindOneById(t *testing.T) {
	const dummyID = 1234
	taskIn := &m.Task{Model: m.Model{ID: dummyID}}

	t.Run("found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindOneById", mock.Anything).Return(taskIn, nil)
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.FindOneById(dummyID)
		assert.Nil(t, err)
		assert.Equal(t, taskIn, taskOut)
	})

	t.Run("not_found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("FindOneById", mock.Anything).Return(taskIn, errors.New(""))
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.FindOneById(dummyID)
		assert.Error(t, err)
		assert.Equal(t, taskIn, taskOut)
	})
}

func TestTaskService_Find(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		tasksIn := []*m.Task{
			{Name: "Test1"},
			{Name: "Test2"},
		}
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Find", mock.Anything).Return(tasksIn, nil)
		taskService := &TaskService{taskStorage: taskStorage}
		tasksOut, err := taskService.Find(make(TaskDemand))
		assert.Nil(t, err)
		assert.Equal(t, tasksIn, tasksOut)
	})

	t.Run("not_found", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Find", mock.Anything).Return([]*m.Task{}, errors.New(""))
		taskService := &TaskService{taskStorage: taskStorage}
		taskOut, err := taskService.Find(make(TaskDemand))
		assert.Error(t, err)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_Update(t *testing.T) {
	var taskIn = &m.Task{Name: "dummy"}

	t.Run("success", func(t *testing.T) {
		var validationErr *v.Errors
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Update", taskIn).Return(taskIn, nil)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{
			taskStorage: taskStorage,
			validator:   validation,
		}
		taskOut, err := taskService.Update(taskIn)

		assert.NotNil(t, taskOut)
		assert.Nil(t, err)
	})

	t.Run("validation_error", func(t *testing.T) {
		validationErr := v.NewErrors()
		validationErr.Add(v.Error{Field: "dummy", Message: "test"})

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{validator: validation}
		taskOut, err := taskService.Update(taskIn)

		assert.Equal(t, validationErr, err)
		assert.Empty(t, taskOut)
	})

	t.Run("database_error", func(t *testing.T) {
		var validationErr *v.Errors
		dbErr := errors.New("simple error")
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Update", taskIn).Return(&m.Task{}, dbErr)

		validation := new(MockedValidation)
		validation.On("Validate", *taskIn).Return(validationErr)

		taskService := &TaskService{
			validator:   validation,
			taskStorage: taskStorage,
		}

		taskOut, err := taskService.Update(taskIn)

		assert.Equal(t, dbErr, err)
		assert.Empty(t, taskOut)
	})
}

func TestTaskService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Delete", mock.Anything).Return(nil)
		taskService := &TaskService{taskStorage: taskStorage}
		err := taskService.Delete(0)
		assert.Nil(t, err)
	})

	t.Run("database_error", func(t *testing.T) {
		errorIn := errors.New("test")
		taskStorage := new(MockedTaskStorage)
		taskStorage.On("Delete", mock.Anything).Return(errorIn)
		taskService := &TaskService{taskStorage: taskStorage}
		err := taskService.Delete(0)
		assert.Equal(t, errorIn, err)
	})
}
