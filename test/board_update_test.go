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

func TestBoardUpdate_OK(t *testing.T) {
	clearTable(t, "boards")

	var (
		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			board   map[string]interface{}

			expectedDescription = stubs[ID-1].description + " UPDATED"
			expectedName        = stubs[ID-1].name + " UPDATED"
		)

		reqString := fmt.Sprintf(
			`{"name":"%s", "description": "%s"}`,
			expectedName,
			expectedDescription,
		)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/boards/%d", ID), bytes.NewBuffer([]byte(reqString)))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &board)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(expectedName, board["name"])
		assert.Equal(expectedDescription, board["description"])
		assert.Equal(float64(ID), board["id"])

		var (
			uID                  uint
			name, description    string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, name, description, created_at, updated_at from boards where id = $1;`, ID).
			Scan(&uID, &name, &description, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedName, name)
		assert.Equal(expectedDescription, description)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestBoardUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/boards/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestBoardUpdate_ValidationError(t *testing.T) {
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
			req, err := http.NewRequest("PUT", "/api/v1/boards/88", bytes.NewBuffer([]byte(jsonStr)))
			must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/88'")

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

func TestBoardUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "boards")

	const (
		name        = "test name"
		description = "test description"
	)

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description)
	)

	req, err := http.NewRequest("PUT", "/api/v1/boards/99", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
