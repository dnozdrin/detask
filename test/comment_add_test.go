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

func TestCommentAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	const (
		text = "any text here"
		task = 1
	)

	var (
		err     error
		comment map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s","task":%d}`, text, task))
	)

	_ = seedTasks(t)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &comment)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/comments/1", response.Header().Get("Location"))
	assert.Equal(1.0, comment["id"])
	assert.Equal(text, comment["text"])
	assert.Equal(1.0, comment["task"])

	var (
		ID, checkTask        uint
		checkText            string
		createdAt, updatedAt time.Time
	)
	err = a.DB.QueryRow(`select id, text, task, created_at, updated_at from comments where id = 1;`).
		Scan(&ID, &checkText, &checkTask, &createdAt, &updatedAt)
	must(t, err, "testing: failed to make database query on column add test")

	assert.Equal(uint(1), ID)
	assert.Equal(text, checkText)
	assert.Equal(uint(task), checkTask)
	assert.WithinDuration(time.Now(), createdAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), updatedAt, maxTestsRunExpected)
}

func TestCommentAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestCommentAdd_ValidationError(t *testing.T) {
	var (
		err  error
		body map[string]interface{}

		text    = makeStringStub(5001)
		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s", "task":%d}`, text, 1))
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 1)
}

func TestCommentAdd_WrongTask(t *testing.T) {
	clearTables(t, "tasks", "comments")

	const (
		text = "any text here"
		task = 1
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s","task":%d}`, text, task))
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a task with the provided ID was not found", body["error"])
}
