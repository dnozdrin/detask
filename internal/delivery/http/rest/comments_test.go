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

func TestCommentAdd_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	const (
		text = "any text here"
		task = 1
	)

	var (
		err     error
		comment map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s","task":%d}`, text, task))
	)

	_ = seedTasks(t)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")

	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &comment)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusCreated, response.Code)
	assert.Equal("/api/v1/comments/1", response.Header().Get("Location"))
	assert.Equal(1.0, comment["id"])
	assert.Equal(text, comment["text"])
	assert.Equal(1.0, comment["task"])

	var (
		ID, checkTask        uint
		checkText            string
		createdAt, updatedAt time.Time
	)
	err = a.DB.QueryRow(`select id, text, task, created_at, updated_at from comments where id = 1;`).
		Scan(&ID, &checkText, &checkTask, &createdAt, &updatedAt)
	must(t, err, "testing: failed to make database query on column add rest_test")

	assert.Equal(uint(1), ID)
	assert.Equal(text, checkText)
	assert.Equal(uint(task), checkTask)
	assert.WithinDuration(time.Now(), createdAt, maxTestsRunExpected)
	assert.WithinDuration(time.Now(), updatedAt, maxTestsRunExpected)
}

func TestCommentAdd_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestCommentAdd_ValidationError(t *testing.T) {
	var (
		err  error
		body map[string]interface{}

		text    = makeStringStub(5001)
		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s", "task":%d}`, text, 1))
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("validation failed", body["error"])
	assert.NotEmpty(body["errors"])
	assert.Len(body["errors"], 1)
}

func TestCommentAdd_WrongTask(t *testing.T) {
	clearTables(t, "tasks", "comments")

	const (
		text = "any text here"
		task = 1
	)
	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"%s","task":%d}`, text, task))
	)

	req, err := http.NewRequest("POST", "/api/v1/comment", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a POST request to '/api/v1/comment'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Equal("a task with the provided ID was not found", body["error"])
}

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

func TestCommentGet_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		comment map[string]interface{}

		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	for k, stub := range stubs {
		ID := k + 1
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/comments/%d", ID), nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comment)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), comment["id"])
		assert.Equal(stub.text, comment["text"])
		assert.Equal(float64(stub.task), comment["task"])
	}
}

func TestCommentGet_NotFound(t *testing.T) {
	clearTable(t, "comments")
	var (
		err  error
		body map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments/66", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments/66'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestCommentList_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		err      error
		comments []map[string]interface{}

		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &comments)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(comments, len(stubs))
	for k, b := range comments {
		assert.Equal(stubs[k].text, b["text"])
		assert.Equal(float64(stubs[k].task), b["task"])
	}
}

func TestCommentList_NoItems(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		err      error
		comments []map[string]interface{}

		assert = testify.New(t)
	)

	req, err := http.NewRequest("GET", "/api/v1/comments", nil)
	must(t, err, "testing: failed to make a GET request to '/api/v1/comments'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &comments)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
	assert.Len(comments, 0)
}

func TestCommentList_Demand(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")
	var (
		comments []map[string]interface{}
		stubs    = seedComments(t)
		assert   = testify.New(t)
	)
	//
	t.Run("demand_by_task_1", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?task=1", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?task=1'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())
		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, len(stubs))
	})
	t.Run("demand_by_task_2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?task=2", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?task=2'")

		_, err = a.DB.Exec(`insert into boards (name, description) values ('rest_test board 2', 'rest_test description 2');`)
		must(t, err, "testing: failed seed a board for comment demand")
		_, err = a.DB.Exec(`insert into columns (name, board, position) values ('rest_test column 2', 2, 2000);`)
		must(t, err, "testing: failed seed a column for comment demand")
		_, err = a.DB.Exec(`insert into tasks (name, "column", position) values ('rest_test task 2', 2, 0.5);`)
		must(t, err, "testing: failed seed a task for comment demand")
		_, err = a.DB.Exec(`insert into comments (text, task) values ('rest_test comment 2', 2);`)
		must(t, err, "testing: failed seed a comment for comment demand")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comments)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(comments, 1)
	})
	t.Run("demand_invalid_param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/comments?dummy=rest_test", nil)
		must(t, err, "testing: failed to make a GET request to '/api/v1/comments?dummy=rest_test'")

		body := make(map[string]interface{})
		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)

		must(t, err, "testing: failed to unmarshal %v", body)

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("invalid filter params", body["error"])
	})
}

func TestCommentUpdate_OK(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	var (
		assert = testify.New(t)
		stubs  = seedComments(t)
	)

	itemsNum := len(stubs)
	for ID := 1; ID <= len(stubs); ID++ {
		itemsNum--

		var (
			jsonReq []byte
			comment map[string]interface{}

			expectedText = stubs[ID-1].text + " UPDATED"
		)

		reqString := fmt.Sprintf(`{"text":"%s", "task": 1}`, expectedText)
		jsonReq = []byte(reqString)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/comments/%d", ID), bytes.NewBuffer(jsonReq))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/%d'", ID)

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &comment)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(float64(ID), comment["id"])
		assert.Equal(expectedText, comment["text"])
		assert.Equal(float64(1), comment["task"])

		var (
			uID, taskID          uint
			text                 string
			createdAt, updatedAt time.Time
		)
		err = a.DB.QueryRow(`select id, text, task, created_at, updated_at from comments where id = $1;`, ID).
			Scan(&uID, &text, &taskID, &createdAt, &updatedAt)
		must(t, err, "testing: failed to make a query on board update rest_test")

		assert.Equal(uint(ID), uID)
		assert.Equal(expectedText, text)
		assert.Equal(uint(1), taskID)
		assert.NotNil(createdAt)
		assert.NotNil(updatedAt)
		assert.True(updatedAt.After(createdAt))
	}
}

func TestCommentUpdate_BadRequest(t *testing.T) {
	var (
		err  error
		body map[string]string

		assert = testify.New(t)
	)

	req, err := http.NewRequest("PUT", "/api/v1/comments/77", bytes.NewBuffer([]byte(`{"name":,,,}`)))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/77'")

	response := executeRequest(req)
	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.NotEmpty(body["error"])
}

func TestCommentUpdate_ValidationError(t *testing.T) {
	var (
		body map[string]interface{}

		assert = testify.New(t)
	)
	t.Run("too_short_text", func(t *testing.T) {
		jsonStr := []byte(`{"text":"", "task": 1}`)
		req, err := http.NewRequest("PUT", "/api/v1/comments/88", bytes.NewBuffer(jsonStr))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/88'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("validation failed", body["error"])
		assert.NotEmpty(body["errors"])
		assert.Len(body["errors"], 1)
	})

	t.Run("too_long_text", func(t *testing.T) {
		jsonStr := []byte(fmt.Sprintf(`{"text":"%s", "task": 1}`, makeStringStub(5001)))
		req, err := http.NewRequest("PUT", "/api/v1/comments/88", bytes.NewBuffer(jsonStr))
		must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/88'")

		response := executeRequest(req)
		err = json.Unmarshal(response.Body.Bytes(), &body)
		must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

		assert.Equal(http.StatusBadRequest, response.Code)
		assert.Equal("validation failed", body["error"])
		assert.NotEmpty(body["errors"])
		assert.Len(body["errors"], 1)
	})
}

func TestCommentUpdate_RecordNotFound(t *testing.T) {
	clearTable(t, "comments")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(`{"text":"rest_test", "task":1}`)
	)

	req, err := http.NewRequest("PUT", "/api/v1/comments/99", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/99'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusNotFound, response.Code)
	assert.Equal("resource was not found", body["error"])
}

func TestCommentUpdate_WrongTask(t *testing.T) {
	clearTables(t, "boards", "columns", "tasks", "comments")

	var (
		err  error
		body map[string]interface{}

		assert  = testify.New(t)
		jsonStr = []byte(fmt.Sprintf(`{"text":"rest_test", "task":999}`))
	)

	_ = seedComments(t)
	req, err := http.NewRequest("PUT", "/api/v1/comments/2", bytes.NewBuffer(jsonStr))
	must(t, err, "testing: failed to make a PUT request to '/api/v1/comments/2'")
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &body)
	must(t, err, "testing: failed to unmarshal %v", response.Body.Bytes())

	assert.Equal(http.StatusOK, response.Code)
}
