package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func clearTables() {
	clearTable("boards")
}

func clearTable(table string) {
	_, _ = a.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
	_, _ = a.DB.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table))
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
