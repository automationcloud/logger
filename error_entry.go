package logger

import (
	errors "github.com/go-errors/errors"
)

/*
var reporter ErrorReporter

func init() {
	reporter = ErrorReporter{
		serviceContext{"workers-autoscaler", os.Getenv("VERSION")},
		os.Stderr,
	}
}

*/

// https://cloud.google.com/error-reporting/docs/formatting-error-messages
type ErrorEntry struct {
	ServiceContext ServiceContext `json:"serviceContext"`
	Message        string         `json:"message"`
	Context        errorContext   `json:"context"`
	logger         *Logger
}

func (er *ErrorEntry) WithHTTPRequest(d HTTPRequestDetails) *ErrorEntry {
	er.Context.HTTPRequest = d
	return er
}

func (er *ErrorEntry) WithUser(u string) *ErrorEntry {
	er.Context.User = u
	return er
}

func (er *ErrorEntry) Report() {
	er.logger.transport.ReportError(er)
}

type errorContext struct {
	ReportLocation stackFrame         `json:"reportLocation"`
	User           string             `json:"user"`
	HTTPRequest    HTTPRequestDetails `json:"httpRequest"`
}

type HTTPRequestDetails struct {
	Method             string `json:"method"`
	Url                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIP           string `json:"remoteIp"`
}

type stackFrame struct {
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	FunctionName string `json:"functionName"`
}

func newErrorContext(err error) errorContext {
	frames := errors.Wrap(err, 3).StackFrames()
	frame := frames[0]
	return errorContext{
		ReportLocation: stackFrame{
			frame.File,
			frame.LineNumber,
			frame.Name,
		},
	}
}
