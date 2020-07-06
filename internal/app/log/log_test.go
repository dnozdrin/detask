// +build unit

package log

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewCORSLogger(t *testing.T) {
	mocked := &LoggerMock{}
	logger := NewCORSLogger(mocked)
	assert.Equal(t, mocked, logger.logger)
}

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
