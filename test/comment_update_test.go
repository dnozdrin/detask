// +build integrational

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCommentUpdate_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	var (
		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			comment map[string]interface{}

			expectedText = stubs[ID-1].text + " UPDATED"
		)

		jsonStr := fmt.Sprintf(`{"text":"%s", "task": 1}`, expectedText)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/comments/%d", ID), bytes.NewBuffer([]byte(jsonStr)))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comment)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), comment["id"])
		assert.Equal(expectedText, comment["text"])
		assert.Equal(float64(1), comment["task"])

		var (
			uID, taskID          uint
			text                 string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, text, task, created_at, updated_at from comments where id = $1;`, ID).
			Scan(&uID, &text, &taskID, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedText, text)
		assert.Equal(uint(1), taskID)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestCommentUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/comments/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestCommentUpdate_ValidationError(t *testing.T) {
	var (
		body map[string]interface{}

		assert = testify.New(t)
	)

	tests := []struct {
		name      string
		jsonStr   string
		errorsNum int
	}{
		{"long_text", fmt.Sprintf(`{"text":"%s"}`, makeStringStub(5001)), 2},
		{"empty_text", `{"text":""}`, 2},
		{"task_required", fmt.Sprintf(`{"text":"%s"}`, makeStringStub(5000)), 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/api/v1/comments/88", bytes.NewBuffer([]byte(test.jsonStr)))
			must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/88'")

			response := executeRequest(req)

			err = json.Unmarshal(response.Body.Bytes(), &body)
			must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

			assert.Equal(http.StatusBadRequest, response.Code)
			assert.Equal("validation failed", body["error"])
			assert.NotEmpty(body["errors"])
			assert.Len(body["errors"], test.errorsNum)
		})
	}
}

func TestCommentUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "comments")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = `{"text":"test", "task":1}`
	)

	req, err := http.NewRequest("PUT", "/api/v1/comments/99", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestCommentUpdate_WrongTask(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = `{"text":"test", "task":999}`
	)

	_ = seedComments(t)
	req, err := http.NewRequest("PUT", "/api/v1/comments/2", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
}
