package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServices(t *testing.T) {
	validator := new(MockedValidation)
	boardStorage := new(MockedBoardStorage)
	columnStorage := new(MockedColumnStorage)
	taskStorage := new(MockedTaskStorage)
	commentStorage := new(MockedCommentStorage)

	services := NewServices(
		validator,
		boardStorage,
		columnStorage,
		taskStorage,
		commentStorage,
	)

	assert.NotNil(t, services.Board)
	assert.NotNil(t, services.Column)
	assert.NotNil(t, services.Task)
	assert.NotNil(t, services.Comment)
}
