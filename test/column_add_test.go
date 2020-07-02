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

func TestColumnAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "test name"
		board            = 1
		position float64 = 1000
	)

	var (
		err    error
		column map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('test board', 'test description');`)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &column)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/columns/1", response.Header().Get("Location"))
	assert.Equal(name, column["name"])
	assert.Equal(position, column["position"])
	assert.Equal(1.0, column["board"])
	assert.Equal(1.0, column["id"])

	var (
		ID                   uint
		checkName            string
		createdAt, updatedAt time.Time
		checkPosition        float64
	)
	err = a.DB.QueryRow(`select id, name, position, created_at, updated_at from "columns" where id = 1;`).
		Scan(&ID, &checkName, &checkPosition, &createdAt, &updatedAt)
	must(t, err, "testing: failed to make database query on column add test")

	assert.Equal(uint(1), ID)
	assert.Equal(name, checkName)
	assert.Equal(position, checkPosition)
	assert.WithinDuration(time.Now(), createdAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), updatedAt, maxTestsRunExpected)
}

func TestColumnAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestColumnAdd_ValidationError(t *testing.T) {
	clearTable(t, "columns")

	const name = ""
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s"}`, name))
	)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 3)
}

func TestColumnAdd_WrongBoard(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%f}`, name, board, position))
	)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a board with the provided ID was not found", body["error"])
}

func TestColumnAdd_PositionDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%.0f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('test board', 'test description');`)
	_, err = a.DB.Exec(`insert into columns (name, board, position) values ('test name 2', $1, $2);`, board, position)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("this position has been already taken", body["error"])
}

func TestColumnAdd_NameDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%.0f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('test board', 'test description');`)
	_, err = a.DB.Exec(`insert into columns (name, board, position) values ($1, $2, 1001);`, name, board)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("a record with this name already exists", body["error"])
}
