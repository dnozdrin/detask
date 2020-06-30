package services

import (
	"github.com/pkg/errors"
)

var (
	ErrRecordNotFound     = errors.New("record was not found")
	ErrRecordAlreadyExist = errors.New("record already exists")
	ErrPositionDuplicate  = errors.New("this record position is already taken")
)
