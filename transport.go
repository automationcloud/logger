package logger

import (
	"encoding/json"
	"io"
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

type TransportGKE struct {
	logWriter   io.Writer
	errorWriter io.Writer
}

func (lt *TransportGKE) SendLog(le *LogEntry) {
	payload := le.payload
	if le.Message != "" {
		payload["message"] = le.Message
	}

	payload["severity"] = le.Severity
	if le.HTTPRequest != nil {
		payload["httpRequest"] = le.HTTPRequest
	}

	if len(le.Labels) > 0 {
		payload["labels"] = le.Labels
	}

	json.NewEncoder(lt.logWriter).Encode(payload)
}

func (lt *TransportGKE) ReportError(le *ErrorEntry) {
	json.NewEncoder(lt.errorWriter).Encode(le)
}
