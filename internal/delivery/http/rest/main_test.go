// todo: consider replacing maps with concrete models
// todo: consider optimizing seeding
// todo: cover with test boundary cases: validation, etc
package rest_test

import (
	"fmt"
	"github.com/dnozdrin/detask/internal/app"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var a app.App

func TestMain(m *testing.M) {
	a.Initialize(
		app.NewDBConfig(
			"postgres",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_MIGRATION_PATH"),
		),
		app.NewConfig(
			app.Test,
			"stderr",
		),
	)
	code := m.Run()
	os.Exit(code)
}

const maxTestsRunExpected = time.Second * 30

func clearTables(t *testing.T, tables ...string) {
	for _, table := range tables {
		clearTable(t, table)
	}
}

func clearTable(t *testing.T, table string) {
	if _, err := a.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
		t.Fatalf("testing: table clearing failed: %v", err)
	}
	if _, err := a.DB.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
		t.Fatalf("testing: sequence restart failed: %v", err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func makeStringStub(len uint) string {
	b := make([]byte, len)
	for i := range b {
		b[i] = 's'
	}
	return string(b)
}

func countItems(t *testing.T, table string) (num int) {
	if err := a.DB.QueryRow(fmt.Sprintf(`select count(*) from %s`, table)).Scan(&num); err != nil {
		t.Fatalf("testing: items counting failed: %v", err)
	}
	return num
}

func must(t *testing.T, err error, message string, arg ...interface{}) {
	if err != nil {
		t.Fatalf(message+" | error: "+err.Error(), arg...)
	}
}

type boardStub struct {
	name, description string
	timestamp         time.Time
}

func seedBoards(t *testing.T) []boardStub {
	var err error
	timestamp := time.Unix(1589932800, 0)
	boards := []boardStub{
		{"rest_test name 1", "rest_test description 1", timestamp},
		{"rest_test name 2", "rest_test description 2", timestamp},
		{"rest_test name 3", "rest_test description 3", timestamp},
	}
	for _, b := range boards {
		_, err = a.DB.Exec(`
			insert into boards (name, description, created_at, updated_at)
			values ($1, $2, $3, $3);`,
			b.name, b.description, b.timestamp,
		)
		must(t, err, "testing: failed to seed boards")
	}

	return boards
}

type columnStub struct {
	name      string
	board     uint
	position  float64
	timestamp time.Time
}

func seedColumns(t *testing.T) []columnStub {
	var (
		err     error
		boardID uint

		timestamp = time.Unix(1589932800, 0)
	)

	err = a.DB.QueryRow(`
			insert into boards (name, description, created_at, updated_at)
			values ($1, $2, $3, $3)
			returning id;`,
		"rest_test name 1", "rest_test description 1", timestamp,
	).Scan(&boardID)
	must(t, err, "testing: failed to seed board for columns")

	columns := []columnStub{
		{"rest_test name 1", boardID, 1000, timestamp},
		{"rest_test name 2", boardID, 2000, timestamp},
		{"rest_test name 3", boardID, 3000, timestamp},
	}
	for _, c := range columns {
		_, err = a.DB.Exec(`
			insert into columns (name, board, position, created_at, updated_at)
			values ($1, $2, $3, $4, $4);`,
			c.name, c.board, c.position, c.timestamp,
		)
		must(t, err, "testing: failed to seed columns")
	}

	return columns
}

type taskStub struct {
	name, description string
	column            uint
	position          float64
	timestamp         time.Time
}

func seedTasks(t *testing.T) []taskStub {
	var (
		err               error
		boardID, columnID uint

		timestamp = time.Unix(1589932800, 0)
	)

	err = a.DB.QueryRow(`
			insert into boards (name, description, created_at, updated_at)
			values ($1, $2, $3, $3)
			returning id;`,
		"rest_test name 1", "rest_test description 1", timestamp,
	).Scan(&boardID)
	must(t, err, "testing: failed to seed board for tasks")

	err = a.DB.QueryRow(`
			insert into columns (name, board, position, created_at, updated_at)
			values ($1, $2, $3, $4, $4)
			returning id;`,
		"rest_test name 1", boardID, 1000, timestamp,
	).Scan(&columnID)
	must(t, err, "testing: failed to seed column for tasks")

	tasks := []taskStub{
		{"rest_test name 1", "rest_test description 1", columnID, 1000, timestamp},
		{"rest_test name 2", "rest_test description 1", columnID, 2000, timestamp},
		{"rest_test name 3", "rest_test description 1", columnID, 3000, timestamp},
	}
	for _, task := range tasks {
		_, err = a.DB.Exec(`
			insert into tasks (name, description, "column", position, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $5);`,
			task.name, task.description, task.column, task.position, task.timestamp,
		)
		must(t, err, "testing: failed to seed tasks")
	}

	return tasks
}

type commentStub struct {
	text      string
	task      uint
	timestamp time.Time
}

func seedComments(t *testing.T) []commentStub {
	var (
		err                       error
		boardID, columnID, taskID uint

		timestamp = time.Unix(1589932800, 0)
	)

	err = a.DB.QueryRow(`
			insert into boards (name, description, created_at, updated_at)
			values ($1, $2, $3, $3)
			returning id;`,
		"rest_test name 1", "rest_test description 1", timestamp,
	).Scan(&boardID)
	must(t, err, "testing: failed to seed a board for comments")

	err = a.DB.QueryRow(`
			insert into columns (name, board, position, created_at, updated_at)
			values ($1, $2, $3, $4, $4)
			returning id;`,
		"rest_test name 1", boardID, 1000, timestamp,
	).Scan(&columnID)
	must(t, err, "testing: failed to seed a column for comments")

	err = a.DB.QueryRow(`
			insert into tasks (name, description, "column", position, created_at, updated_at)
			values ('rest_test name 1', 'rest_test description 1', $1, 1000, $2, $2)
			returning id;`,
		columnID, timestamp,
	).Scan(&taskID)
	must(t, err, "testing: failed to seed a task for comments")

	comments := []commentStub{
		{"rest_test text 1", taskID, timestamp},
		{"rest_test text 2", taskID, timestamp},
		{"rest_test text 3", taskID, timestamp},
	}
	for _, c := range comments {
		_, err = a.DB.Exec(`
			insert into comments (text, task, created_at, updated_at)
			values ($1, $2, $3, $3);`,
			c.text, c.task, c.timestamp,
		)
		must(t, err, "testing: failed to seed comments")
	}

	return comments
}
