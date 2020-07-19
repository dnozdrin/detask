// +build integrational

package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/health")
	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"status\":\"OK\"}", response.Body.String())
}
