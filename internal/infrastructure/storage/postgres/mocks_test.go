// +build unit

package postgres

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Errorf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Error(args ...interface{}) {
	l.Called(args)
}

func (l *LoggerMock) Fatalf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Fatal(args ...interface{}) {
	l.Called(args)
}

func (l *LoggerMock) Infof(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Info(args ...interface{}) {
	l.Called(args)
}

func (l *LoggerMock) Warnf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Warn(args ...interface{}) {
	l.Called(args)
}

func (l *LoggerMock) Debugf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Debug(args ...interface{}) {
	l.Called(args)
}

type DBMock struct {
	mock.Mock
}

func (db *DBMock) Begin() (*sql.Tx, error) {
	returnValues := db.Called()
	return returnValues.Get(0).(*sql.Tx), returnValues.Error(1)
}

func (db *DBMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	returnValues := db.Called(query, args)
	return returnValues.Get(0).(*sql.Rows), returnValues.Error(1)
}

func (db *DBMock) QueryRow(query string, args ...interface{}) *sql.Row {
	returnValues := db.Called(query, args)
	return returnValues.Get(0).(*sql.Row)
}

func (db *DBMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	returnValues := db.Called(query, args)
	return returnValues.Get(0).(sql.Result), returnValues.Error(1)
}

func (db *DBMock) Prepare(query string) (*sql.Stmt, error) {
	returnValues := db.Called(query)
	return returnValues.Get(0).(*sql.Stmt), returnValues.Error(1)
}
