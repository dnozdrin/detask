package test

import (
	"fmt"
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBoardDelete(t *testing.T) {
	clearTable(t, "boards")

	var (
		num int

		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/boards/%d", ID), nil)
		must(t, err, "testing: failed to make a DELETE request to '/api/v1/boards/%d'", ID)

		response := executeRequest(req)
		num = countItems(t, "boards")

		assert.Equal(http.StatusNoContent, response.Code)
		assert.Empty(response.Body.Bytes())
		assert.Equal(itemsNum, num)
	}
}
