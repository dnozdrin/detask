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
