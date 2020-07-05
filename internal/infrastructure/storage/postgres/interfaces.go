package postgres

import "database/sql"

type db interface {
	querier
	txBeginner
}

type querier interface {
	preparator
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type txBeginner interface {
	Begin() (*sql.Tx, error)
}

type preparator interface {
	Prepare(query string) (*sql.Stmt, error)
}
