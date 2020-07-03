package test

import (
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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
