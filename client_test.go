package logging

import (
	"testing"
)

type mockTransport struct {
	logEntry   *LogEntry
	errorEntry *ErrorEntry
}

func (t *mockTransport) SendLog(le *LogEntry) error {
	t.logEntry = le
	return nil
}

func (t *mockTransport) ReportError(er *ErrorEntry) error {
	t.errorEntry = er
	return nil
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
