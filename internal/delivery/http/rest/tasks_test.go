package rest_test

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
		name                = "rest_test name"
		description         = "rest_test description"
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
	must(t, err, "testing: failed to make database query on column add rest_test")

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
		name             = "rest_test name"
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
		name             = "rest_test name"
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

	_, err = a.DB.Exec(`insert into boards (name, description) values ($1, 'rest_test description');`, name)
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

func TestTaskDelete(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	var (
		assert   = testify.New(t)
		stubs    = seedTasks(t)
		itemsNum = len(stubs)
		ID       = 1
	)

	for ; ID <= len(stubs); ID++ {
		itemsNum--
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/tasks/%d", ID), nil)
		must(t, err, "testing: failed to make a DELETE request to '/api/v1/tasks/%d'", ID)
		response := executeRequest(req)

		assert.Equal(http.StatusNoContent, response.Code)
		assert.Empty(response.Body.Bytes())
		assert.Equal(itemsNum, countItems(t, "tasks"))
	}
}

func TestTaskGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")
	var (
		task map[string]interface{}

		assert = testify.New(t)
		stubs  = seedTasks(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &task)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), task["id"])
		assert.Equal(stub.name, task["name"])
		assert.Equal(stub.description, task["description"])
		assert.Equal(stub.position, task["position"])
		assert.Equal(float64(stub.column), task["column"])
	}
}

func TestTaskGet_NotFound(t *testing.T) {
	clearTables(t, "tasks")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/tasks/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/tasks/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

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

		_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board 2', 'rest_test description 2');`)
		must(t, err, "testing: failed seed a board for task demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('rest_test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for task demand")
		_, err = a.DB.Exec(`insert into tasks (name, "column", position) values ('rest_test task N', 2, 0.5);`)
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
		req, err := http.NewRequest("GET", "/api/v1/tasks?dummy=rest_test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/tasks?dummy=rest_test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)
		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}

func TestTaskUpdate_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks")

	var (
		assert = testify.New(t)
		stubs  = seedTasks(t)
	)

	_, err := a.DB.Exec(`insert into columns (name, board, position) values ('rest_test column 2', 1, 2000);`)
	must(t, err, "testing: failed seed a column for task update")

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			jsonReq []byte
			task    map[string]interface{}

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
		jsonReq = []byte(reqString)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/tasks/%d", ID), bytes.NewBuffer(jsonReq))
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
		must(t, err, "testing: failed to make a query on board update rest_test")

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
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"description":"%s", "position": 10}`, makeStringStub(5001)))
	)

	req, err := http.NewRequest("PUT", "/api/v1/columns/88", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/88'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 2)
}

func TestTaskUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "tasks")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(`{"name":"rest_test", "description":"rest_test", "position": 1, "column": 1}`)
	)

	req, err := http.NewRequest("PUT", "/api/v1/tasks/99", bytes.NewBuffer(jsonStr))
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
		jsonStr = []byte(fmt.Sprintf(`{"name":"rest_test", "description":"rest_test", "position": 1000, "column": 1}`))
	)

	_ = seedTasks(t)
	req, err := http.NewRequest("PUT", "/api/v1/tasks/2", bytes.NewBuffer(jsonStr))
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
		jsonStr = []byte(fmt.Sprintf(`{"name":"rest_test", "description":"rest_test", "position": 1000, "column": 999}`))
	)

	_ = seedTasks(t)
	req, err := http.NewRequest("PUT", "/api/v1/tasks/2", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/tasks/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a column with the provided ID was not found", body["error"])
}
