// +build integrational

package test

import (
	"encoding/json"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTaskList_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")
	var (
		err   error
		tasks []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedTasks(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/tasks", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/tasks'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &tasks)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(tasks, len(stubs))
	for k, b := range tasks {
		assert.Equal(stubs[k].name, b["name"])
		assert.Equal(stubs[k].description, b["description"])
		assert.Equal(float64(stubs[k].column), b["column"])
		assert.Equal(stubs[k].position, b["position"])
	}
}

func TestTaskList_NoItems(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")
	var (
		err   error
		tasks []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/tasks", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/tasks'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &tasks)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(tasks, 0)
}

func TestTaskList_Demand(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")
	var (
		tasks  []map[string]interface{}
		stubs  = seedTasks(t)
		assert = testify.New(t)
	)

	t.Run("demand_by_column_1", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/tasks?column=1", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?column=1'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &tasks)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(tasks, len(stubs))
	})
	t.Run("demand_by_column_2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/tasks?column=2", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?column=2'")

		_, err = a.DB.Exec(`insert into boards (name, description) values ('test board 2', 'test description 2');`)
		must(t, err, "testing: failed seed a board for task demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for task demand")
		_, err = a.DB.Exec(`insert into tasks (name, "column", position) values ('test task N', 2, 0.5);`)
		must(t, err, "testing: failed seed a task for task demand")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &tasks)

		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(tasks, 1)
	})
	t.Run("demand_by_board_1", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/tasks?board=1", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?board=1'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &tasks)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(tasks, len(stubs))
	})
	t.Run("demand_by_board_2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/tasks?board=2", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?board=2'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &tasks)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(tasks, 1)
	})
	t.Run("demand_invalid_param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/tasks?dummy=test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?dummy=test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)
		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}
