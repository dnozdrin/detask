package services

import (
	"github.com/pkg/errors"
)

var (
	ErrRecordNotFound     = errors.New("record was not found")
	ErrRecordAlreadyExist = errors.New("record already exists")

	ErrNameDuplicate     = errors.New("a record with this name already exists")
	ErrPositionDuplicate = errors.New("this position has been already taken")

	ErrBoardRelation  = errors.New("a board with the provided ID was not found")
	ErrColumnRelation = errors.New("a column with the provided ID was not found")
	ErrTaskRelation   = errors.New("a task with the provided ID was not found")
)
