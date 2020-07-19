package services

import (
	"github.com/pkg/errors"
)

var (
	// ErrRecordNotFound is used for cases when the requested record was not found
	ErrRecordNotFound = errors.New("record was not found")

	// ErrRecordAlreadyExist is used for cases when there is an attempt to create a record
	// which already exists (ID constraints violation).
	ErrRecordAlreadyExist = errors.New("record already exists")

	// ErrNameDuplicate is used for cases when there is an attempt to create or modify a record
	// and the new name violates unique constraints.
	ErrNameDuplicate = errors.New("a record with this name already exists")

	// ErrPositionDuplicate is used for cases when there is an attempt to create or modify a record
	// and the new position violates unique constraints.
	ErrPositionDuplicate = errors.New("this position has been already taken")

	// ErrBoardRelation is used for cases when there is an attempt to create a relation with a
	// board that does not exist in the system.
	ErrBoardRelation = errors.New("a board with the provided ID was not found")

	// ErrColumnRelation is used for cases when there is an attempt to create a relation with a
	// column that does not exist in the system.
	ErrColumnRelation = errors.New("a column with the provided ID was not found")

	// ErrTaskRelation is used for cases when there is an attempt to create a relation with a
	// task that does not exist in the system.
	ErrTaskRelation = errors.New("a task with the provided ID was not found")

	// ErrLastColumn is used for cases when there is an attempt to delete the last column on a board.
	ErrLastColumn = errors.New("the last column can not be deleted")

	// ErrTargetColumn is used for cases when the target column for tasks on a column deletion was not found
	ErrTargetColumn = errors.Errorf("columns storage: target column for tasks transfer not found")
)
