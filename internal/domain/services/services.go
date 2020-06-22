package services

import (
	"github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/pkg/errors"
)

type Services struct {
	Board   *BoardService
	Column  *ColumnService
	Task    *TaskService
	Comment *CommentService
}

var (
	ErrRecordNotFound     = errors.New("record was not found")
	ErrRecordAlreadyExist = errors.New("record already exists")
)

func NewServices(
	validator validation.Validator,
	boardStorage BoardStorage,
	columnStorage ColumnStorage,
	taskStorage TaskStorage,
	commentStorage CommentStorage,
) *Services {
	return &Services{
		Board:   NewBoardService(validator, boardStorage),
		Column:  NewColumnService(validator, columnStorage),
		Task:    NewTaskService(validator, taskStorage),
		Comment: NewCommentService(validator, commentStorage),
	}
}
