package logger

import (
	"testing"
)

type mockTransport struct {
	logEntry   *LogEntry
	errorEntry *ErrorEntry
}

func (t *mockTransport) SendLog(le *LogEntry) {
	t.logEntry = le
}

func (t *mockTransport) ReportError(er *ErrorEntry) {
	t.errorEntry = er
}

var client *Client

func init() {
	DefaultTransport = &mockTransport{}
	client = NewClient("test", "0.0.0")
}

func TestNewLogger(t *testing.T) {
	e := client.NewLogEntry()
	if e.Message != "" {
		t.Error("expect blank message")
	}

	if e.client.Transport != DefaultTransport {
		t.Error("expect default log transport to be used")
	}
}
