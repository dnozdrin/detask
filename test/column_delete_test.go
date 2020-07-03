package test

import (
	"encoding/json"
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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
