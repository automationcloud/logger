package logger

import (
	"os"
)

var DefaultTransport Transport

func init() {
	DefaultTransport = &TransportGKE{os.Stdout, os.Stderr}
}

type Transport interface {
	LogSender
	ErrorReporter
}

type LogSender interface {
	SendLog(*LogEntry)
}

type ErrorReporter interface {
	ReportError(*ErrorEntry)
}
