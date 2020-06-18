package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

type TaskStorage interface {
	Save(task *m.Task) (*m.Task, error)
	FindById(id uint) (*m.Task, error)
	FindAll() ([]*m.Task, error)
	Update(task *m.Task) (*m.Task, error)
	Delete(id uint) error
}

type TaskService struct {
	taskStorage    TaskStorage
	commentStorage CommentStorage
	validator      v.Validator
}

// NewTaskService is a task service constructor
func NewTaskService(taskStorage TaskStorage, validator v.Validator) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
		validator:   validator,
	}
}

// Create will create a new task with the provided payload. Returns the
// operation result with possible validation or saving errors
func (t *TaskService) Create(task *m.Task) (*m.Task, v.Result) {
	var result v.Result
	if result = t.validator.Validate(*task); !result.IsValid() {
		return nil, result
	}

	task, result.Error = t.taskStorage.Save(task)

	return task, result
}

// FindAll will return all not deleted comments and an error in case
// it occurred while fetching records from the storage
func (t *TaskService) FindAll() ([]*m.Task, error) {
	return t.taskStorage.FindAll()
}

// FindOneById will return a pointer to the comment requested by id and
// an error in case it occurred while fetching the record from the storage
func (t *TaskService) FindOneById(ID uint) (*m.Task, error) {
	return t.taskStorage.FindById(ID)
}

// Update will update the task record. Returns the operation result
// with possible validation or saving errors
func (t *TaskService) Update(task *m.Task) (*m.Task, v.Result) {
	var result v.Result
	if result = t.validator.Validate(*task); !result.IsValid() {
		return nil, result
	}

	task, result.Error = t.taskStorage.Update(task)

	return task, result
}

// Delete will delete a record with the given ID
func (t *TaskService) Delete(ID uint) error {
	return t.taskStorage.Delete(ID)
}
