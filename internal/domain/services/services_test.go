package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServices(t *testing.T) {
	bs := new(MockedBoardStorage)
	cls := new(MockedColumnStorage)
	ts := new(MockedTaskStorage)
	cmts := new(MockedCommentStorage)
	v := new(MockedValidation)

	services := NewServices(v, bs, cls, ts, cmts)

	assert.NotNil(t, services.Board)
	assert.NotNil(t, services.Column)
	assert.NotNil(t, services.Task)
	assert.NotNil(t, services.Comment)
}
