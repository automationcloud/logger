package logging

import (
	"net/http"
	"runtime/debug"

	errors "github.com/go-errors/errors"
)

type ErrorEntry struct {
	Error        error
	Request      *http.Request
	User         string
	Stack        string
	CodeLocation StackFrame
	client       *Client
}

type StackFrame struct {
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	FunctionName string `json:"functionName"`
}

func (l *Client) NewErrorEntry(err error) *ErrorEntry {
	return &ErrorEntry{
		Error:        err,
		Stack:        string(debug.Stack()),
		CodeLocation: captureLocation(err, 3),
		client:       l,
	}
}

func (er *ErrorEntry) WithRequest(r *http.Request) *ErrorEntry {
	er.Request = r
	return er
}

func (er *ErrorEntry) WithUser(u string) *ErrorEntry {
	er.User = u
	return er
}

func (er *ErrorEntry) Report() {
	er.client.Transport.ReportError(er)
}

func captureLocation(err error, skip int) StackFrame {
	frames := errors.Wrap(err, skip).StackFrames()
	frame := frames[0]
	return StackFrame{
		frame.File,
		frame.LineNumber,
		frame.Name,
	}
}
