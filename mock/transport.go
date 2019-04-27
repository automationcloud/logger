package mock

import "github.com/automationcloud/logging"

type Transport struct {
	ErrorEntry *logging.ErrorEntry
	LogEntry   *logging.LogEntry
}

func (ft *Transport) SendLog(l *logging.LogEntry) {
	ft.LogEntry = l
}

func (ft *Transport) ReportError(l *logging.ErrorEntry) {
	ft.ErrorEntry = l
}
