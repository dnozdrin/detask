package log

// CORSLogger is a logger wrapper for CORS package
type CORSLogger struct {
	logger Logger
}

// NewCORSLogger is a CORSLogger constructor
func NewCORSLogger(logger Logger) CORSLogger {
	return CORSLogger{logger}
}

// Printf uses fmt.Sprintf to log a templated message.
func (l CORSLogger) Printf(format string, args ...interface{}) {
	l.logger.Debugf(format, args)
}
