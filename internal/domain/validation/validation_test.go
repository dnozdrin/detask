package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErrors(t *testing.T) {
	e := NewErrors()
	assert.Equal(t, e.msg, failMessage)
	assert.Empty(t, e.errors)
	assert.IsType(t, []Error{}, e.errors)
}

func TestErrorsMessage(t *testing.T) {
	e := NewErrors()
	e.msg = "test 123"
	assert.Equal(t, e.msg, e.Error())
}

func TestErrorsNum(t *testing.T) {
	e := NewErrors()
	assert.Equal(t, e.Num(), 0)
	e.errors = append(e.errors, Error{})
	assert.Equal(t, e.Num(), 1)
	e.errors = append(e.errors, Error{})
	assert.Equal(t, e.Num(), 2)
}

func TestErrorsAdd(t *testing.T) {
	e := NewErrors()
	assert.Len(t, e.errors, 0)
	e.Add(Error{})
	assert.Equal(t, e.Num(), 1)
	e.Add(Error{})
	assert.Equal(t, e.Num(), 2)
}

func TestMarshalJSON(t *testing.T) {
	e := NewErrors()
	e.errors = append(e.errors, Error{Field: "name", Message: "incorrect"})
	result, err := e.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte("{\"error\":\"validation failed\",\"errors\":[{\"field\":\"name\",\"message\":\"incorrect\"}]}"), result)
}