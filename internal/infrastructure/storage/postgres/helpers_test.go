// +build unit

package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeferred(t *testing.T) {
	logger := new(LoggerMock)
	logger.On("Errorf", "%v", mock.Anything).Return().Once()

	tests := []struct {
		name string
		f    func() error
	}{
		{"sql_error", func() error { return sql.ErrTxDone }},
		{"any_other_error", func() error { return errors.New("test") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deferred(logger, tt.f)
		})
	}
}
