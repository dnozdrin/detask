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

func TestColumnAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "rest_test name"
		board            = 1
		position float64 = 1000
	)

	var (
		err    error
		column map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board', 'rest_test description');`)

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
	must(t, err, "testing: failed to make database query on column add rest_test")

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
		name             = "rest_test name"
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
		name             = "rest_test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%.0f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board', 'rest_test description');`)
	_, err = a.DB.Exec(`insert into columns (name, board, position) values ('rest_test name 2', $1, $2);`, board, position)

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
		name             = "rest_test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s","board":%d,"position":%.0f}`, name, board, position))
	)

	_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board', 'rest_test description');`)
	_, err = a.DB.Exec(`insert into columns (name, board, position) values ($1, $2, 1001);`, name, board)

	req, err := http.NewRequest("POST", "/api/v1/column", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/column'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("a record with this name already exists", body["error"])
}

func TestColumnDelete(t *testing.T) {
	clearTables(t, "boards", "columns")

	var (
		assert   = testify.New(t)
		stubs    = seedColumns(t)
		itemsNum = len(stubs)
		ID       = 1
		body     map[string]interface{}
	)

	t.Run("delete_ok", func(t *testing.T) {
		for ; ID <= len(stubs)-1; ID++ {
			itemsNum--
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/columns/%d", ID), nil)
			must(t, err, "testing: failed to make a DELETE request to '/api/v1/columns/%d'", ID)
			response := executeRequest(req)

			assert.Equal(http.StatusNoContent, response.Code)
			assert.Empty(response.Body.Bytes())
			assert.Equal(itemsNum, countItems(t, "columns"))
		}
	})
	t.Run("delete_last_column", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/columns/%d", ID), nil)
		must(t, err, "testing: failed to make a DELETE request to '/api/v1/columns/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("the last column on the board can not be deleted", body["error"])
		assert.Equal(1, countItems(t, "columns"))
	})
}

func TestColumnGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns")
	var (
		column map[string]interface{}

		assert = testify.New(t)
		stubs  = seedColumns(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/columns/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &column)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), column["id"])
		assert.Equal(stub.name, column["name"])
		assert.Equal(stub.position, column["position"])
		assert.Equal(float64(stub.board), column["board"])
	}
}

func TestColumnGet_NotFound(t *testing.T) {
	clearTables(t, "columns")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/columns/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/columns/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestColumnList_OK(t *testing.T) {
	clearTables(t, "boards", "columns")
	var (
		err     error
		columns []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedColumns(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/columns", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/columns'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &columns)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(columns, len(stubs))
	for k, b := range columns {
		assert.Equal(stubs[k].name, b["name"])
		assert.Equal(stubs[k].position, b["position"])
		assert.Equal(float64(stubs[k].board), b["board"])
	}
}

func TestColumnList_NoItems(t *testing.T) {
	clearTable(t, "columns")
	var (
		err     error
		columns []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/columns", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/columns'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &columns)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(columns, 0)
}

func TestColumnList_Demand(t *testing.T) {
	clearTables(t, "boards", "columns")
	var (
		comments []map[string]interface{}
		stubs    = seedColumns(t)
		assert   = testify.New(t)
	)

	t.Run("demand_by_board_1", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/columns?board=1", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns?board=1'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())
		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, len(stubs))
	})
	t.Run("demand_by_board_2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/columns?board=2", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns?board=2'")

		_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board 2', 'rest_test description 2');`)
		must(t, err, "testing: failed seed a board for comment demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('rest_test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for comment demand")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, 1)
	})
	t.Run("demand_invalid_param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/columns?dummy=rest_test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns?dummy=rest_test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)

		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}

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
			jsonReq []byte
			column  map[string]interface{}

			expectedName     = stubs[ID-1].name + " UPDATED"
			expectedPosition = stubs[ID-1].position / 2
		)

		reqString := fmt.Sprintf(
			`{"name":"%s", "position": %f, "board": %d}`,
			expectedName,
			expectedPosition,
			stubs[ID-1].board,
		)
		jsonReq = []byte(reqString)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/columns/%d", ID), bytes.NewBuffer(jsonReq))
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
		must(t, err, "testing: failed to make a query on board update rest_test")

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
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":""}`))
	)

	req, err := http.NewRequest("PUT", "/api/v1/columns/88", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/88'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 3)
}

func TestColumnUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "columns")

	const name = "rest_test name"

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "position": %d, "board": %d}`, name, 1, 1))
	)

	req, err := http.NewRequest("PUT", "/api/v1/columns/99", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestColumnUpdate_PositionDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "rest_test name"
		board            = 1
		position float64 = 1000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "board":%d, "position":%f}`, name, board, position))
	)

	_ = seedColumns(t)
	req, err := http.NewRequest("PUT", "/api/v1/columns/2", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("this position has been already taken", body["error"])
}

func TestColumnUpdate_NameDuplicate(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name             = "rest_test name 1"
		board            = 1
		position float64 = 2000
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "board":%d, "position":%f}`, name, board, position))
	)

	_ = seedColumns(t)

	req, err := http.NewRequest("PUT", "/api/v1/columns/2", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/columns/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusConflict, response.Code)
	assert.Equal("a record with this name already exists", body["error"])
}
