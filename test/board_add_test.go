package test

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
	clearTable("boards")
	clearTable("columns")

	const (
		name        = "test name 1"
		description = "test description"
	)

	var (
		err error

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description))
	)

	req, _ := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/boards/1", response.Header().Get("Location"))

	var board map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &board)

	assert.Nil(err)
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
	assert.Nil(err)

	assert.Equal(uint(1), bID)
	assert.Equal(name, bName)
	assert.Equal(description, bDescription)
	assert.NotNil(bCreatedAt)
	assert.NotNil(bUpdatedAt)

	assert.Equal(uint(1), cID)
	assert.Equal("Default", cName)
	assert.Equal(float64(1000), cPosition)
	assert.NotNil(cCreatedAt)
	assert.NotNil(cUpdatedAt)
}

func TestBoardAdd_BadRequest(t *testing.T) {
	var (
		err error

		assert  = testify.New(t)
		jsonStr = []byte(`{"name":,,,}`)
	)

	req, _ := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)
	assert.Equal(http.StatusBadRequest, response.Code)

	var body map[string]string
	err = json.Unmarshal(response.Body.Bytes(), &body)

	assert.Nil(err)
	assert.NotEmpty(body["error"])
}

func TestBoardAdd_ValidationError(t *testing.T) {
	clearTable("boards")

	const name = ""
	var (
		err error

		description = makeStringStub(1001)
		assert      = testify.New(t)
		jsonStr     = []byte(fmt.Sprintf(`{"name":"%s", "description": "%s"}`, name, description))
	)

	_, _ = a.DB.Exec(`insert into boards (name, description) values ($1, $2);`, name, description)
	req, err := http.NewRequest("POST", "/api/v1/board", bytes.NewBuffer(jsonStr))
	assert.Nil(err)
	response := executeRequest(req)

	assert.Equal(http.StatusBadRequest, response.Code)

	var body map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Nil(err)
	assert.NotEmpty(body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 2)
}
