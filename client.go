package logging

import (
	"context"
)

// Key to use when setting the request ID.
type ctxKeyLogger int

// RequestIDKey is the key that holds the unique request ID in a request context.
const LoggerKey ctxKeyLogger = 0

// Client is an assembly point for transport and service providing logging and
// error reporting capabilities.
type Client struct {
	ServiceContext
	Transport
}

// NewClient creates a new client with specified service name and version.
func NewClient(service, version string) *Client {
	return &Client{
		ServiceContext: ServiceContext{
			Service: service,
			Version: version,
		},
		Transport: DefaultTransport,
	}
}

// FromContext restores a client from a context.
func FromContext(ctx context.Context) *Client {
	l, ok := ctx.Value(LoggerKey).(Client)
	if ok {
		return &l
	}

	return NewClient("", "")
}

// ToContext puts client to a context.
func (l *Client) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, LoggerKey, *l)
}

// WrapLogSender adds a middleware to log processing stack, returning a new client.
//
// Example use case for it is adding request-specific logging information to a
// logging client before passing it to a business layer. This way logs will contain
// necessary logging information without breaking boundaries between transport
// and business layer.
func (l *Client) WrapLogSender(fn func(*LogEntry)) *Client {
	return &Client{
		ServiceContext: l.ServiceContext,
		Transport: transportMiddleware{
			originalTransport: l.Transport,
			logSender:         fn,
		},
	}
}

// WrapErrorReporter adds a middleware to error processing stack, returning a new client.
//
// Similarly to WrapLogSender, WrapErrorReporter helps to keep business layer
// free of transport layer information by wrapping Transport.
//
// Use-case: call ErrorEntry.WithRequest(r) to attach request info to error
// reported from a business layer.
func (l *Client) WrapErrorReporter(fn func(*ErrorEntry)) *Client {
	return &Client{
		ServiceContext: l.ServiceContext,
		Transport: transportMiddleware{
			originalTransport: l.Transport,
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
