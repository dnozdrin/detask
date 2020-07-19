// +build unit

package rest

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIDVarError_Columns(t *testing.T) {
	logger := new(LoggerMock)
	logger.On("Errorf", mock.Anything, mock.Anything).Return()

	router := new(RouteAwareMock)
	router.On("GetIDVar", new(http.Request)).Return(uint(1), errors.New("test error"))

	boardHandler := ColumnHandler{log: logger, router: router, resp: &responder{log: logger}}

	tests := []struct {
		name   string
		method func(http.ResponseWriter, *http.Request)
	}{
		{name: "GetOneById", method: boardHandler.GetOneById},
		{name: "Update", method: boardHandler.Update},
		{name: "Delete", method: boardHandler.Delete},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(boardHandler.GetOneById)
			handler.ServeHTTP(recorder, &http.Request{})

			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		})
	}
}
