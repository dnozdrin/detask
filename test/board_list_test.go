package test

import (
	"encoding/json"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type boardStub struct {
	name, description string
}

func TestBoardList_OK(t *testing.T) {
	clearTable("boards")
	assert := testify.New(t)

	stubs, err := seedBoards()
	assert.Nil(err)

	req, _ := http.NewRequest("GET", "/api/v1/boards", nil)
	response := executeRequest(req)
	assert.Equal(http.StatusOK, response.Code)

	var boards []map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &boards)
	assert.Nil(err)

	assert.Len(boards, len(stubs))

	for k, b := range boards {
		assert.Equal(stubs[k].name, b["name"])
		assert.Equal(stubs[k].description, b["description"])
	}
}

func TestBoardList_NoItems(t *testing.T) {
	clearTable("boards")
	assert := testify.New(t)

	req, _ := http.NewRequest("GET", "/api/v1/boards", nil)
	response := executeRequest(req)
	assert.Equal(http.StatusOK, response.Code)

	var boards []map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &boards)
	assert.Nil(err)
	assert.Len(boards, 0)
}

func seedBoards() ([]boardStub, error) {
	var err error
	boards := []boardStub{
		{"test name 1", "test description 1"},
		{"test name 2", "test description 2"},
		{"test name 3", "test description 3"},
		{"test name 4", "test description 4"},
	}
	for _, b := range boards {
		_, err = a.DB.Exec(`insert into boards (name, description) values ($1, $2);`, b.name, b.description)
		if err != nil {
			return boards, err
		}
	}

	return boards, nil
}
