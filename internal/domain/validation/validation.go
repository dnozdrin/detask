package validation

import (
	"encoding/json"
)

const failMessage = "validation failed"

// Error is a container for a validation error description
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Errors is a container for set of validation errors
type Errors struct {
	msg    string
	errors []Error
}

// NewErrors is a validation errors container constructor
func NewErrors() *Errors {
	return &Errors{
		msg:    failMessage,
		errors: make([]Error, 0),
	}
}

// Error will return errors container general message
func (e *Errors) Error() string {
	return e.msg
}

// Num will return the number of errors in the errors container
func (e *Errors) Num() int {
	return len(e.errors)
}

// Add will add the provided error to the errors container
func (e *Errors) Add(err Error) {
	e.errors = append(e.errors, err)
}

// MarshalJSON provides correct marshaling for Errors type
func (e *Errors) MarshalJSON() ([]byte, error) {
	input := struct {
		Msg    string  `json:"error"`
		Errors []Error `json:"errors"`
	}{e.msg, e.errors}
	return json.Marshal(input)
}

// Validator represents an interface for a struct validation
type Validator interface {
	// Validate should return a pointer to Errors. In case there are no validation errors,
	// the referenced Errors with have nil value
	Validate(interface{}) *Errors
}
