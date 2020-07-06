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

func TestTaskUpdate_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	var (
		assert = testify.New(t)
		stubs  = seedTasks(t)
	)

	_, err := a.DB.Exec(`insert into columns (name, board, position) values ('test column 2', 1, 2000);`)
	must(t, err, "testing: failed seed a column for task update")

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			task map[string]interface{}

			expectedName        = stubs[ID-1].name + " UPDATED"
			expectedColumn      = 2
			expectedDescription = stubs[ID-1].description + " UPDATED"
			expectedPosition    = stubs[ID-1].position / 2
		)

		reqString := fmt.Sprintf(
			`{"name":"%s", "description":"%s", "position": %f, "column": %d}`,
			expectedName,
			expectedDescription,
			expectedPosition,
			expectedColumn,
		)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/tasks/%d", ID), bytes.NewBuffer([]byte(reqString)))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &task)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), task["id"])
		assert.Equal(expectedName, task["name"])
		assert.Equal(expectedDescription, task["description"])
		assert.Equal(expectedPosition, task["position"])
		assert.Equal(float64(expectedColumn), task["column"])

		var (
			uID, columnID        uint
			position             float64
			name, description    string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, name, description, position, "column", created_at, updated_at from tasks where id = $1;`, ID).
			Scan(&uID, &name, &description, &position, &columnID, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedName, name)
		assert.Equal(expectedDescription, description)
		assert.Equal(expectedPosition, position)
		assert.Equal(uint(expectedColumn), columnID)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestTaskUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/tasks/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestTaskUpdate_ValidationError(t *testing.T) {
	var (
		body map[string]interface{}

		assert = testify.New(t)
	)

	tests := []struct {
		name      string
		jsonStr   string
		errorsNum int
	}{
		{"long_description", fmt.Sprintf(`{"name":"test", "description":"%s"}`, makeStringStub(5001)), 3},
		{"empty_name_empty_description", `{"name":""}`, 4},
		{"empty_name_with_description", fmt.Sprintf(`{"name":"", "description":"%s"}`, makeStringStub(5000)), 3},
		{"name_set_empty_description", `{"name":"test"}`, 3},
		{"position_required", `{"name":"test", "column": 1}`, 2},
		{"column_required", `{"name":"test", "position": 1000}`, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/api/v1/tasks/88", bytes.NewBuffer([]byte(test.jsonStr)))
			must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/88'")

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

func TestTaskUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "tasks")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = `{"name":"test", "description":"test", "position": 1, "column": 1}`
	)

	req, err := http.NewRequest("PUT", "/api/v1/tasks/99", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestTaskUpdate_PositionDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = fmt.Sprintf(`{"name":"test", "description":"test", "position": 1000, "column": 1}`)
	)

	_ = seedTasks(t)
	req, err := http.NewRequest("PUT", "/api/v1/tasks/2", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("this position has been already taken", body["error"])
}

func TestTaskUpdate_WrongColumn(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = fmt.Sprintf(`{"name":"test", "description":"test", "position": 1000, "column": 999}`)
	)

	_ = seedTasks(t)
	req, err := http.NewRequest("PUT", "/api/v1/tasks/2", bytes.NewBuffer([]byte(jsonStr)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a column with the provided ID was not found", body["error"])
}
