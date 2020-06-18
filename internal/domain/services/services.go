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

func NewServices(v validation.Validator, bs BoardStorage, cls ColumnStorage, ts TaskStorage, cmts CommentStorage) *Services {
	return &Services{
		Board:   NewBoardService(bs, v),
		Column:  NewColumnService(cls, v),
		Task:    NewTaskService(ts, v),
		Comment: NewCommentService(cmts, v),
	}
}
