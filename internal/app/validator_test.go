// +build unit

package app

import (
	"testing"

	in "github.com/dnozdrin/detask/internal/domain/validation"
	validate "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewValidator(t *testing.T) {
	validator := validate.New()
	logger := new(LoggerMock)
	v := NewValidator(validator, logger)

	assert.Equal(t, logger, v.log)
	assert.Equal(t, validator, v.validate)
}

func TestSuccessValidation(t *testing.T) {
	tests := []struct {
		name   string
		target interface{}
	}{
		{
			name: "anonymous_1",
			target: struct {
				Test1 string `validate:"min=1"`
				Test2 string `validate:"max=2"`
			}{"1", "11"},
		},
		{
			name: "anonymous_2",
			target: struct {
				Test string `validate:"required"`
			}{":-)"},
		},
	}
	for _, test := range tests {
		validator := Validator{log: new(LoggerMock), validate: validate.New()}
		t.Run("validate_"+test.name, func(t *testing.T) {
			err := validator.Validate(test.target)
			assert.IsType(t, &in.Errors{}, err)
			assert.Nil(t, err)
		})
	}
}

func TestInvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		target interface{}
	}{
		{name: "map[string]string", target: make(map[string]string)},
		{name: "[]int", target: make([]int, 0)},
		{name: "string", target: mock.Anything},
		{name: "int", target: 41},
		{name: "float64", target: 4.1},
		{name: "bool", target: true},
	}

	for _, test := range tests {
		t.Run("validate_invalid_type_"+test.name, func(t *testing.T) {
			logger := new(LoggerMock)
			logger.On("Errorf", "non-struct value passed for validation: %s", []interface{}{test.name})
			validator := Validator{log: logger, validate: validate.New()}
			err := validator.Validate(test.target)
			assert.IsType(t, &in.Errors{}, err)
			assert.EqualError(t, err, "validation failed")
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		target    interface{}
		errorsNum int
	}{
		{
			name: "anonymous_1",
			target: struct {
				Test1 string `validate:"min=1"`
				Test2 string `validate:"max=2"`
			}{"", "1111"},
			errorsNum: 2,
		},
		{
			name: "anonymous_2",
			target: struct {
				Test string `validate:"required"`
			}{""},
			errorsNum: 1,
		},
		{
			name: "anonymous_3",
			target: struct {
				// not supported
				Test string `validate:"gt=0"`
			}{""},
			errorsNum: 1,
		},
	}
	for _, test := range tests {
		validator := Validator{log: new(LoggerMock), validate: validate.New()}
		t.Run("empty_"+test.name, func(t *testing.T) {
			err := validator.Validate(test.target)
			assert.IsType(t, &in.Errors{}, err)
			assert.Equal(t, test.errorsNum, err.Num())
		})
	}
}
