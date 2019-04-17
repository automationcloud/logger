package logging

// DefaultTransport used in NewClient. Initialised as TransportGKE.
var DefaultTransport Transport

func init() {
	DefaultTransport = NewTransportGKE()
}

// Transport is an interface for sending logs and reporting errors
type Transport interface {
	LogSender
	ErrorReporter
}

// LogSender is an interface for sending logs
type LogSender interface {
	SendLog(*LogEntry)
}

// ErrorReporter is an interface for reporting errors
type ErrorReporter interface {
	ReportError(*ErrorEntry)
}
