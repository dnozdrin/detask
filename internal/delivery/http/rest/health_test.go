// +build unit

package rest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck_Status(t *testing.T) {
	loggerMock := new(LoggerMock)
	loggerMock.On("Info", mock.Anything).Return()
	healthCheck := NewHealthCheck(loggerMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheck.Status)
	handler.ServeHTTP(recorder, &http.Request{})

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"status":"OK"}`, recorder.Body.String())
}
