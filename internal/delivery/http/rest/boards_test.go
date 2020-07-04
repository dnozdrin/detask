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

func TestBoardAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns")

	const (
		name        = "rest_test name"
		description = "rest_test description"
	)

	var (
		err   error
		board map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "description":"%s"}`, name, description))
	)

	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/board'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &board)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/boards/1", response.Header().Get("Location"))
	assert.Equal(name, board["name"])
	assert.Equal(description, board["description"])
	assert.Equal(1.0, board["id"])

	var (
		bID, cID                                       uint
		bName, bDescription, cName                     string
		bCreatedAt, bUpdatedAt, cCreatedAt, cUpdatedAt time.Time
		cPosition                                      float64
	)
	err = a.DB.QueryRow(
		`select
				b.id, b.name, b.description, b.created_at, b.updated_at,
       			c.id, c.name, c.position, c.created_at, c.updated_at
       		from boards b join columns c on b.id = c.board where b.id = 1;`,
	).Scan(
		&bID, &bName, &bDescription, &bCreatedAt, &bUpdatedAt,
		&cID, &cName, &cPosition, &cCreatedAt, &cUpdatedAt,
	)
	must(t, err, "testing: failed to make database query on board add rest_test")

	assert.Equal(uint(1), bID)
	assert.Equal(name, bName)
	assert.Equal(description, bDescription)
	assert.WithinDuration(time.Now(), bCreatedAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), bUpdatedAt, maxTestsRunExpected)

	assert.Equal(uint(1), cID)
	assert.Equal("Default", cName)
	assert.Equal(float64(1000), cPosition)
	assert.WithinDuration(time.Now(), cCreatedAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), cUpdatedAt, maxTestsRunExpected)
}

func TestBoardAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/board'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestBoardAdd_ValidationError(t *testing.T) {
	clearTable(t, "boards")

	const name = ""
	var (
		err  error
		body map[string]interface{}

		description = makeStringStub(1001)
		assert      = testify.New(t)
		jsonStr     = []byte(fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description))
	)

	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/board'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 2)
}

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

func TestBoardGet_OK(t *testing.T) {
	clearTable(t, "boards")
	var (
		board map[string]interface{}

		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/boards/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/boards/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &board)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), board["id"])
		assert.Equal(stub.name, board["name"])
		assert.Equal(stub.description, board["description"])
	}
}

func TestBoardGet_NotFound(t *testing.T) {
	clearTable(t, "boards")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestBoardList_OK(t *testing.T) {
	clearTable(t, "boards")
	var (
		err    error
		boards []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &boards)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(boards, len(stubs))
	for k, b := range boards {
		assert.Equal(stubs[k].name, b["name"])
		assert.Equal(stubs[k].description, b["description"])
	}
}

func TestBoardList_NoItems(t *testing.T) {
	clearTable(t, "boards")
	var (
		err    error
		boards []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/boards", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/boards'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &boards)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(boards, 0)
}

func TestBoardUpdate_OK(t *testing.T) {
	clearTable(t, "boards")

	var (
		assert = testify.New(t)
		stubs  = seedBoards(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			jsonReq []byte
			board   map[string]interface{}

			expectedDescription = stubs[ID-1].description + " UPDATED"
			expectedName        = stubs[ID-1].name + " UPDATED"
		)

		reqString := fmt.Sprintf(
			`{"name":"%s", "description": "%s"}`,
			expectedName,
			expectedDescription,
		)
		jsonReq = []byte(reqString)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/boards/%d", ID), bytes.NewBuffer(jsonReq))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &board)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(expectedName, board["name"])
		assert.Equal(expectedDescription, board["description"])
		assert.Equal(float64(ID), board["id"])

		var (
			uID                  uint
			name, description    string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, name, description, created_at, updated_at from boards where id = $1;`, ID).
			Scan(&uID, &name, &description, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update rest_test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedName, name)
		assert.Equal(expectedDescription, description)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestBoardUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/boards/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestBoardUpdate_ValidationError(t *testing.T) {
	clearTable(t, "boards")

	const name = ""
	var (
		err  error
		body map[string]interface{}

		description = makeStringStub(1001)
		assert      = testify.New(t)
		jsonStr     = []byte(fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description))
	)

	req, err := http.NewRequest("PUT", "/api/v1/boards/88", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/88'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 2)
}

func TestBoardUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "boards")

	const (
		name        = "rest_test name"
		description = "rest_test description"
	)

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description))
	)

	req, err := http.NewRequest("PUT", "/api/v1/boards/99", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/boards/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}
