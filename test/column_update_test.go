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

func TestColumnUpdate_OK(t *testing.T) {
	clearTables(t, "boards", "columns")

	var (
		assert = testify.New(t)
		stubs  = seedColumns(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			column  map[string]interface{}

			expectedName     = stubs[ID-1].name + " UPDATED"
			expectedPosition = stubs[ID-1].position / 2
		)

		jsonStr := fmt.Sprintf(
			`{"name":"%s", "position": %f, "board": %d}`,
			expectedName,
			expectedPosition,
			stubs[ID-1].board,
		)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/columns/%d", ID), bytes.NewBuffer([]byte(jsonStr)))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &column)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(expectedName, column["name"])
		assert.Equal(expectedPosition, column["position"])
		assert.Equal(float64(ID), column["id"])

		var (
			uID, board           uint
			position             float64
			name                 string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, name, position, board, created_at, updated_at from columns where id = $1;`, ID).
			Scan(&uID, &name, &position, &board, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedName, name)
		assert.Equal(expectedPosition, position)
		assert.Equal(stubs[ID-1].board, board)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestColumnUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/columns/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestColumnUpdate_ValidationError(t *testing.T) {
	var (
		body map[string]interface{}

		assert = testify.New(t)
	)

	tests := []struct {
		name      string
		jsonStr   string
		errorsNum int
	}{
		{"long_name", fmt.Sprintf(`{"name":"%s"}`, makeStringStub(256)), 3},
		{"empty_name", `{"name":""}`, 3},
		{"name_set", `{"name":"test"}`, 2},
		{"position_required", `{"name":"test", "board": 1}`, 1},
		{"board_required", `{"name":"test", "position": 1000}`, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/api/v1/columns/88", bytes.NewBuffer([]byte(test.jsonStr)))
			must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/88'")

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

func TestColumnUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "columns")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = fmt.Sprintf(`{"name":"%s", "position": %d, "board": %d}`, "test name", 1, 1)
	)

	req, err := http.NewRequest("PUT", "/api/v1/columns/99", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestColumnUpdate_PositionDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = `{"name":"test name", "board":1, "position":1000}`
	)

	_ = seedColumns(t)
	req, err := http.NewRequest("PUT", "/api/v1/columns/2", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("this position has been already taken", body["error"])
}

func TestColumnUpdate_NameDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = `{"name":"test name 1", "board":1, "position":1000}`
	)

	_ = seedColumns(t)

	req, err := http.NewRequest("PUT", "/api/v1/columns/2", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("a record with this name already exists", body["error"])
}
