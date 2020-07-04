// +build integrational

package test

import (
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTaskGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")
	var (
		task map[string]interface{}

		assert = testify.New(t)
		stubs  = seedTasks(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &task)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), task["id"])
		assert.Equal(stub.name, task["name"])
		assert.Equal(stub.description, task["description"])
		assert.Equal(stub.position, task["position"])
		assert.Equal(float64(stub.column), task["column"])
	}
}

func TestTaskGet_NotFound(t *testing.T) {
	clearTables(t, "tasks")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/tasks/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/tasks/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
