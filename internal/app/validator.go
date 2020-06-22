package app

import (
	"reflect"
	"strings"

	in "github.com/dnozdrin/detask/internal/domain/validation"
	pkg "github.com/go-playground/validator/v10"
)

// Validator represents an implementation of validation feature
type Validator struct {
	log      Logger
	validate *pkg.Validate
}

// NewValidator is a Validator constructor
func NewValidator(validate *pkg.Validate, log Logger) *Validator {
	return &Validator{
		validate: validate,
		log:      log,
	}
}

// Validate returns a pointer to Errors. In case there are no validation errors,
// the referenced Errors with have nil value
func (v Validator) Validate(target interface{}) *in.Errors {
	var result *in.Errors

	reflected := reflect.ValueOf(target)
	if reflected.Kind() != reflect.Struct {
		result := in.NewErrors()
		result.Add(in.Error{Message: "invalid input dataset"})
		v.log.Errorf("non-struct value passed for validation: %s", reflected.Type().String())

		return result
	}

	err := v.validate.Struct(target)
	if err != nil {
		return convertValidationErrors(target, err)
	}

	return result
}

func convertValidationErrors(target interface{}, err error) *in.Errors {
	result := in.NewErrors()
	validationErrors := err.(pkg.ValidationErrors)

	reflected := reflect.ValueOf(target)
	for _, e := range validationErrors {
		field, _ := reflected.Type().FieldByName(e.StructField())

		var name string
		if name = field.Tag.Get("json"); name == "" {
			name = strings.ToLower(e.StructField())
		}

		result.Add(in.Error{Field: name, Message: formatMessage(e, name)})
	}

	return result
}

func formatMessage(err pkg.FieldError, name string) (message string) {
	switch err.Tag() {
	case "required":
		message = name + " is required"
	case "max":
		message = name + " must be of " + err.Param() + " symbols max"
	case "min":
		message = name + " must be of " + err.Param() + " symbols min"
	default:
		message = name + " is invalid"
	}

	return message
}
