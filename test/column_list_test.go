package test

import (
	"encoding/json"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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

		_, err = a.DB.Exec(`insert into boards (name, description) values ('test board 2', 'test description 2');`)
		must(t, err, "testing: failed seed a board for comment demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for comment demand")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, 1)
	})
	t.Run("demand_invalid_param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/columns?dummy=test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/columns?dummy=test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)

		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}
