package validation

import "github.com/pkg/errors"

var (
	ErrValidationFailed = errors.New("validation failed")
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Errors []Error

type Result struct {
	Error    error  `json:"-"`
	ErrorMsg string `json:"error"`
	Errors   Errors `json:"errors"`
}

func NewResult(err error) Result {
	var message string
	if err != nil {
		message = err.Error()
	}
	return Result{
		Error:    err,
		ErrorMsg: message,
		Errors:   make(Errors, 0),
	}
}

func (v Result) IsValid() bool {
	return len(v.Errors) == 0
}

type Validator interface {
	Validate(interface{}) Result
}
