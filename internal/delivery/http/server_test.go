// +build unit

package http

import (
	testify "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewServer(t *testing.T) {
	var (
		assert  = testify.New(t)
		handler = new(HandlerMock)
		logger  = new(LoggerMock)
		server  = NewServer(handler, logger)
	)

	assert.Equal(logger, server.log)
	assert.Equal(handler, server.http.Handler)
	assert.Equal(writeTimeout, server.http.WriteTimeout)
	assert.Equal(readTimeout, server.http.ReadTimeout)
	assert.IsType(&http.Server{}, server.http)
}

func Test_server_Start(t *testing.T) {
	const addr = ":24080"
	var (
		serverErr error
		assert    = testify.New(t)
		handler   = new(HandlerMock)
		logger    = new(LoggerMock)
		server    = NewServer(handler, logger)

		serverStarted = make(chan struct{})
		serverDone    = make(chan struct{})
	)

	logger.On("Infof", "starting http server on address %s", []interface{}{addr}).
		Return().
		Twice()

	go func() {
		close(serverStarted)
		serverErr = server.Start(addr)
		defer close(serverDone)
	}()

	<-serverStarted
	server.shutdown()
	<-serverDone

	assert.Nil(serverErr)
	logger.AssertNumberOfCalls(t, "Infof", 1)
	logger.AssertNotCalled(t, "Errorf", 0)
}
