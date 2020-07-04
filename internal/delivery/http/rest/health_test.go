// +build unit

package rest

import (
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

	if recorder.Code != http.StatusOK {
		t.Errorf("bad response code, wanted %v got %v",
			recorder.Code, http.StatusOK)
	}
	expectedBody := `{"status":"OK"}`
	if recorder.Body.String() != expectedBody {
		t.Errorf("bad response body, wanted %v got %v", recorder.Body, expectedBody)
	}
}
