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

func TestBoardAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns")

	var (
		err   error
		board map[string]interface{}

		assert      = testify.New(t)
		name        = makeStringStub(500)
		description = makeStringStub(1000)
		jsonStr     = fmt.Sprintf(`{"name":"%s", "description":"%s"}`, name, description)
	)

	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/board'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &board)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/boards/1", response.Header().Get("Location"))
	assert.Equal(name, board["name"])
	assert.Equal(description, board["description"])
	assert.Equal(1.0, board["id"])

	var (
		bID, cID                                       uint
		bName, bDescription, cName                     string
		bCreatedAt, bUpdatedAt, cCreatedAt, cUpdatedAt time.Time
		cPosition                                      float64
	)
	err = a.DB.QueryRow(
		`select
				b.id, b.name, b.description, b.created_at, b.updated_at,
       			c.id, c.name, c.position, c.created_at, c.updated_at
       		from boards b join columns c on b.id = c.board where b.id = 1;`,
	).Scan(
		&bID, &bName, &bDescription, &bCreatedAt, &bUpdatedAt,
		&cID, &cName, &cPosition, &cCreatedAt, &cUpdatedAt,
	)
	must(t, err, "testing: failed to make database query on board add test")

	assert.Equal(uint(1), bID)
	assert.Equal(name, bName)
	assert.Equal(description, bDescription)
	assert.WithinDuration(time.Now(), bCreatedAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), bUpdatedAt, maxTestsRunExpected)

	assert.Equal(uint(1), cID)
	assert.Equal("Default", cName)
	assert.Equal(float64(1000), cPosition)
	assert.WithinDuration(time.Now(), cCreatedAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), cUpdatedAt, maxTestsRunExpected)
}

func TestBoardAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/board'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestBoardAdd_ValidationError(t *testing.T) {
	var (
		body    map[string]interface{}
		jsonStr string

		assert          = testify.New(t)
		longDescription = makeStringStub(1001)
		longName        = makeStringStub(501)
	)

	type board struct {
		name, description string
	}

	tests := []struct {
		name      string
		board     board
		errorsNum int
	}{
		{"long_description", board{"test", longDescription}, 1},
		{"long_name", board{longName, "test"}, 1},
		{"long_description_empty_name", board{"", longDescription}, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonStr = fmt.Sprintf(`{"name":"%s", "description": "%s"}`, test.board.name, test.board.description)
			req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer([]byte(jsonStr)))
			must(t, err, "testing: failed to make a POST request to '/api/v1/board'")
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
