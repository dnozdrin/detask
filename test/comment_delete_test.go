// +build integrational

package test

import (
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCommentDelete(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	var (
		assert   = testify.New(t)
		stubs    = seedComments(t)
		itemsNum = len(stubs)
		ID       = 1
	)

	for ; ID <= len(stubs); ID++ {
		itemsNum--
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/comments/%d", ID), nil)
		must(t, err, "testing: failed to make a DELETE request to '/api/v1/comments/%d'", ID)
		response := executeRequest(req)

		assert.Equal(http.StatusNoContent, response.Code)
		assert.Empty(response.Body.Bytes())
		assert.Equal(itemsNum, countItems(t, "comments"))
	}
}
