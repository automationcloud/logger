package mock

import "github.com/automationcloud/logging"

type Transport struct {
	errorEntry *logging.ErrorEntry
	logEntry   *logging.LogEntry
}

func (ft *Transport) SendLog(l *logging.LogEntry) {
	ft.logEntry = l
}

func (ft *Transport) ReportError(l *logging.ErrorEntry) {
	ft.errorEntry = l
}
