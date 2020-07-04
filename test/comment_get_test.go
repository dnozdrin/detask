// +build integrational

package test

import (
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCommentGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		comment map[string]interface{}

		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/comments/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comment)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), comment["id"])
		assert.Equal(stub.text, comment["text"])
		assert.Equal(float64(stub.task), comment["task"])
	}
}

func TestCommentGet_NotFound(t *testing.T) {
	clearTable(t, "comments")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
