package test

import (
	"encoding/json"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCommentList_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		err      error
		comments []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &comments)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(comments, len(stubs))
	for k, b := range comments {
		assert.Equal(stubs[k].text, b["text"])
		assert.Equal(float64(stubs[k].task), b["task"])
	}
}

func TestCommentList_NoItems(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		err      error
		comments []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &comments)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(comments, 0)
}

func TestCommentList_Demand(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		comments []map[string]interface{}
		stubs    = seedComments(t)
		assert   = testify.New(t)
	)
	//
	t.Run("demand_by_task_1", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?task=1", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?task=1'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())
		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, len(stubs))
	})
	t.Run("demand_by_task_2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?task=2", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?task=2'")

		_, err = a.DB.Exec(`insert into boards (name, description) values ('test board 2', 'test description 2');`)
		must(t, err, "testing: failed seed a board for comment demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for comment demand")
		_, err = a.DB.Exec(`insert into tasks (name, "column", position) values ('test task 2', 2, 0.5);`)
		must(t, err, "testing: failed seed a task for comment demand")
		_, err = a.DB.Exec(`insert into comments (text, task) values ('test comment 2', 2);`)
		must(t, err, "testing: failed seed a comment for comment demand")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, 1)
	})
	t.Run("demand_invalid_param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?dummy=test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?dummy=test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)

		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}
