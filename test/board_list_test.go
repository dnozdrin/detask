// +build integrational

package test

import (
	"encoding/json"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBoardList_OK(t *testing.T) {
	clearTable(t, "boards")
	var (
		err    error
		boards []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &boards)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(boards, len(stubs))
	for k, b := range boards {
		assert.Equal(stubs[k].name, b["name"])
		assert.Equal(stubs[k].description, b["description"])
	}
}

func TestBoardList_NoItems(t *testing.T) {
	clearTable(t, "boards")
	var (
		err    error
		boards []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &boards)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(boards, 0)
}
