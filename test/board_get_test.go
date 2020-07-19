// +build integrational

package test

import (
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBoardGet_OK(t *testing.T) {
	clearTable(t, "boards")
	var (
		board map[string]interface{}

		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/boards/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/boards/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &board)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), board["id"])
		assert.Equal(stub.name, board["name"])
		assert.Equal(stub.description, board["description"])
	}
}

func TestBoardGet_NotFound(t *testing.T) {
	clearTable(t, "boards")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
