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
	b := make([]byte,len)
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
		must(t, err, "testing: while seeding boards")
	}

	return boards
}
