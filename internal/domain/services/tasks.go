package services

import (
	m "github.com/dnozdrin/detask/internal/domain/models"
	v "github.com/dnozdrin/detask/internal/domain/validation"
)

// TaskService is an interactor for work with tasks
type TaskService struct {
	validator   v.Validator
	taskStorage TaskStorage
}

// NewTaskService is a task service constructor
func NewTaskService(validator v.Validator, taskStorage TaskStorage) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
		validator:   validator,
	}
}

// Create will create a new task with the provided payload. Returns the
// operation result with possible validation or saving errors
func (t *TaskService) Create(task *m.Task) (*m.Task, error) {
	if err := t.validator.Validate(*task); err != nil {
		return nil, err
	}

	return t.taskStorage.Save(task)
}

// Find will return all tasks that meet the provided demand and an
// error in case it occurred while fetching records from the storage
func (t *TaskService) Find(demand TaskDemand) ([]*m.Task, error) {
	return t.taskStorage.Find(demand)
}

// FindOneById will return a pointer to the task requested by id and
// an error in case it occurred while fetching the record from the storage
func (t *TaskService) FindOneById(ID uint) (*m.Task, error) {
	return t.taskStorage.FindOneById(ID)
}

// Update will update the task record. Returns the operation result
// with possible validation or saving errors
func (t *TaskService) Update(task *m.Task) (*m.Task, error) {
	if err := t.validator.Validate(*task); err != nil {
		return nil, err
	}

	return t.taskStorage.Update(task)
}

// Delete will delete a record with the given ID
func (t *TaskService) Delete(ID uint) error {
	return t.taskStorage.Delete(ID)
}
