package http

import (
	"context"
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

func Test_server_StartFail(t *testing.T) {
	const addr = ":8080"
	var (
		assert  = testify.New(t)
		handler = new(HandlerMock)
		logger  = new(LoggerMock)
		server  = NewServer(handler, logger)

		serverStarted = make(chan struct{})
		serverDone    = make(chan struct{})
	)

	logger.On("Infof", "starting http server on address %s", []interface{}{addr}).Return()

	fakeServer := &http.Server{Addr: addr, Handler: handler}
	go func() {
		serverStarted <- struct{}{}
		if err := fakeServer.ListenAndServe(); err != nil {
			t.Fatalf("failed to start fake server: %v", err)
		}
	}()

	go func() {
		<-serverStarted
		err := server.Start(addr)
		assert.Error(err)
		serverDone <- struct{}{}
	}()

	<-serverDone
	if err := fakeServer.Shutdown(context.Background()); err != nil {
		t.Fatalf("failed to start fake server: %v", err)
	}

	close(serverStarted)
	close(serverDone)
}
