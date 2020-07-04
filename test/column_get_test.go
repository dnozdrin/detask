// +build integrational

package test

import (
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestColumnGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns")
	var (
		column map[string]interface{}

		assert = testify.New(t)
		stubs  = seedColumns(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/columns/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &column)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), column["id"])
		assert.Equal(stub.name, column["name"])
		assert.Equal(stub.position, column["position"])
		assert.Equal(float64(stub.board), column["board"])
	}
}

func TestColumnGet_NotFound(t *testing.T) {
	clearTables(t, "columns")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/columns/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/columns/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
