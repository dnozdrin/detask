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
	tests := []struct {
		name   string
		errors []Error
		json []byte
	}{
		{
			name: "no_errors",
			errors: []Error{},
			json: []byte("{\"error\":\"validation failed\",\"errors\":[]}"),
		},
		{
			name: "1_error",
			errors: []Error{{Field: "name", Message: "incorrect"}},
			json: []byte("{\"error\":\"validation failed\",\"errors\":[{\"field\":\"name\",\"message\":\"incorrect\"}]}"),
		},
		{
			name: "2_errors",
			errors: []Error{{Field: "dummy", Message: "error"}, {Field: "tests", Message: "test"}},
			json: []byte("{\"error\":\"validation failed\",\"errors\":[{\"field\":\"dummy\",\"message\":\"error\"},{\"field\":\"tests\",\"message\":\"test\"}]}"),
		},
	}
	for _, test := range tests {
		t.Run("validation_"+test.name, func(t *testing.T) {
			e := NewErrors()
			e.errors = test.errors
			result, err := e.MarshalJSON()
			assert.Nil(t, err)
			assert.Equal(t, test.json, result)
		})
	}
}
