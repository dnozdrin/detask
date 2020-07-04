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

func TestTaskAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	const (
		name                = "test name"
		description         = "test description"
		column              = 1
		position    float64 = 1000
	)

	var (
		err  error
		task map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(
			`{"name":"%s","description":"%s","column":%d,"position":%f}`,
			name,
			description,
			column,
			position,
		),
		)
	)

	_ = seedColumns(t)

	req, err := http.NewRequest("POST", "/api/v1/task", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/task'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &task)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/tasks/1", response.Header().Get("Location"))
	assert.Equal(1.0, task["id"])
	assert.Equal(name, task["name"])
	assert.Equal(1.0, task["column"])
	assert.Equal(position, task["position"])

	var (
		ID, checkColumn      uint
		checkName            string
		checkDescription     string
		createdAt, updatedAt time.Time
		checkPosition        float64
	)
	err = a.DB.QueryRow(`select id, name, description, "column", position, created_at, updated_at from tasks where id = 1;`).
		Scan(&ID, &checkName, &checkDescription, &checkColumn, &checkPosition, &createdAt, &updatedAt)
	must(t, err, "testing: failed to make database query on column add test")

	assert.Equal(uint(1), ID)
	assert.Equal(name, checkName)
	assert.Equal(description, checkDescription)
	assert.Equal(uint(column), checkColumn)
	assert.Equal(position, checkPosition)
	assert.WithinDuration(time.Now(), createdAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), updatedAt, maxTestsRunExpected)
}

func TestTaskAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/task", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/task'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestTaskAdd_ValidationError(t *testing.T) {
	const name = ""
	var (
		err  error
		body map[string]interface{}

		description = makeStringStub(5001)
		assert      = testify.New(t)
		jsonStr     = []byte(fmt.Sprintf(`{"name":"%s", "description":"%s"}`, name, description))
	)

	req, err := http.NewRequest("POST", "/api/v1/task", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/task'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 4)
}

func TestTaskAdd_WrongColumn(t *testing.T) {
	clearTables(t, "columns", "tasks")

	const (
		name             = "test name"
		column           = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","column":%d,"position":%f}`, name, column, position))
	)

	req, err := http.NewRequest("POST", "/api/v1/task", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/task'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a column with the provided ID was not found", body["error"])
}

func TestTaskAdd_PositionDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	const (
		name             = "test name"
		board            = 1
		column           = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","column":%d,"position":%f}`, name, column, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ($1, 'test description');`, name)
	_, err = a.DB.Exec(`insert into columns (name, board, position) values ($1, $2, $3);`, name, board, position)
	_, err = a.DB.Exec(`insert into tasks (name, "column", position) values ($1, $2, $3);`, name, column, position)

	req, err := http.NewRequest("POST", "/api/v1/task", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/task'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("this position has been already taken", body["error"])
}
