// +build unit

package rest

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
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

type RouteAwareMock struct {
	mock.Mock
}

func (raw *RouteAwareMock) GetURL(name string, params ...string) (*url.URL, error) {
	returnValues := raw.Called(name, params)
	return returnValues.Get(0).(*url.URL), returnValues.Error(1)
}

func (raw *RouteAwareMock) GetIDVar(r *http.Request) (uint, error) {
	returnValues := raw.Called(r)
	return returnValues.Get(0).(uint), returnValues.Error(1)
}
