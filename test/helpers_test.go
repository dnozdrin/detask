package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
		t.Fatalf(message, arg...)
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
		{"test name 1", "test description 1", timestamp},
		{"test name 2", "test description 2", timestamp},
		{"test name 3", "test description 3", timestamp},
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
	name string
	board uint
	position float64
	timestamp         time.Time
}


func seedColumns(t *testing.T) []columnStub {
	var (
		err error
		boardID uint

		timestamp = time.Unix(1589932800, 0)
	)

	err = a.DB.QueryRow(`
			insert into boards (name, description, created_at, updated_at)
			values ($1, $2, $3, $3)
			returning id;`,
		"test name 1", "test description 1", timestamp,
	).Scan(&boardID)
	must(t, err, "testing: failed to seed board for columns")

	columns := []columnStub{
		{"test name 1", boardID, 1000, timestamp},
		{"test name 2", boardID, 2000, timestamp},
		{"test name 3", boardID, 3000, timestamp},
	}
	for _, b := range columns {
		_, err = a.DB.Exec(`
			insert into columns (name, board, position, created_at, updated_at)
			values ($1, $2, $3, $4, $4);`,
			b.name, b.board, b.position, b.timestamp,
		)
		must(t, err, "testing: failed to seed columns")
	}

	return columns
}