package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"status\":\"OK\"}", response.Body.String())
}
