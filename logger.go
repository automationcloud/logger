package logger

import (
	"context"
	"runtime/debug"
)

// Key to use when setting the request ID.
type ctxKeyLogger int

// RequestIDKey is the key that holds the unique request ID in a request context.
const LoggerKey ctxKeyLogger = 0

type Logger struct {
	ServiceContext ServiceContext
	transport      Transport
}

func NewLogger(service, version string) *Logger {
	return &Logger{
		ServiceContext: ServiceContext{
			Service: service,
			Version: version,
		},
		transport: DefaultTransport,
	}
}

func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(LoggerKey).(Logger)
	if ok {
		return &l
	}

	return NewLogger("", "")
}

func (l *Logger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, LoggerKey, *l)
}

func (l *Logger) NewLogEntry() *LogEntry {
	return &LogEntry{
		Labels:  make(map[string]string),
		Payload: make(map[string]interface{}),
		logger:  l,
	}
}

func (l *Logger) NewError(err error) *ErrorEntry {
	return &ErrorEntry{
		ServiceContext: l.ServiceContext,
		Message:        err.Error() + "\n" + string(debug.Stack()),
		Context:        newErrorContext(err),
		logger:         l,
	}
}

func (l *Logger) WrapLogSender(fn func(*LogEntry)) *Logger {
	return &Logger{
		ServiceContext: l.ServiceContext,
		transport: transportMiddleware{
			originalTransport: l.transport,
			logSender:         fn,
		},
	}
}

func (l *Logger) WrapErrorReporter(fn func(*ErrorEntry)) *Logger {
	return &Logger{
		ServiceContext: l.ServiceContext,
		transport: transportMiddleware{
			originalTransport: l.transport,
			errorReporter:     fn,
		},
	}
}

type transportMiddleware struct {
	originalTransport Transport
	logSender         func(*LogEntry)
	errorReporter     func(*ErrorEntry)
}

func (tm transportMiddleware) SendLog(l *LogEntry) {
	if tm.logSender != nil {
		tm.logSender(l)
	}
	tm.originalTransport.SendLog(l)
}

func (tm transportMiddleware) ReportError(e *ErrorEntry) {
	if tm.errorReporter != nil {
		tm.errorReporter(e)
	}
	tm.originalTransport.ReportError(e)
}
