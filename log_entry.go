package logging

import (
	"net/http"
	"time"
)

type Severity int

const (
	DEFAULT   Severity = 0   // The log entry has no assigned severity level.
	DEBUG     Severity = 100 // Debug or trace information.
	INFO      Severity = 200 // Routine information, such as ongoing status or performance.
	NOTICE    Severity = 300 // Normal but significant events, such as start up, shut down, or a configuration change.
	WARNING   Severity = 400 // Warning events might cause problems.
	ERROR     Severity = 500 // Error events are likely to cause problems.
	CRITICAL  Severity = 600 // Critical events cause more severe problems or outages.
	ALERT     Severity = 700 // A person must take an action immediately.
	EMERGENCY Severity = 800 // One or more systems are unusable.
)

// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type LogEntry struct {
	Severity    Severity          `json:"severity"`
	Message     string            `json:"message,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	HTTPRequest *HTTPRequest      `json:"httpRequest,omitempty"`
	Request     *http.Request
	Payload     map[string]interface{}
	client      *Client
}

// https://godoc.org/cloud.google.com/go/logging#HTTPRequest
type HTTPRequest struct {
	Request *http.Request

	// The size of the HTTP request message in bytes, including the request headers and the request body.
	RequestSize int64

	// The response code indicating the status of response. Examples: 200, 404.
	Status int

	// The size of the HTTP response message sent back to the client, in bytes, including the response headers and the response body.
	ResponseSize int64
	RemoteIP     string
	ServerIP     string
	Latency      time.Duration

	CacheLookup                    bool
	CacheHit                       bool
	CacheValidatedWithOriginServer bool
	CacheFillBytes                 string
}

// NewLogEntry creates log entry ready for sending.
func (l *Client) NewLogEntry() *LogEntry {
	return &LogEntry{
		Labels:  make(map[string]string),
		Payload: make(map[string]interface{}),
		client:  l,
	}
}

// WithPayload adds data by key to structured payload.
func (le *LogEntry) WithPayload(key string, data interface{}) *LogEntry {
	le.Payload[key] = data
	return le
}

// WithRequest adds http roundtrip information to structured payload.
func (le *LogEntry) WithRequest(hr *HTTPRequest) *LogEntry {
	le.HTTPRequest = hr
	return le
}

// Log sends an log entry to a transport with a severity and message.
func (le *LogEntry) Log(s Severity, msg string) {
	le.Message = msg
	le.Severity = s
	le.client.Transport.SendLog(le)
}

func (le *LogEntry) Debug(msg string) {
	le.Log(DEBUG, msg)
}

func (le *LogEntry) Info(msg string) {
	le.Log(INFO, msg)
}

func (le *LogEntry) Warn(msg string) {
	le.Log(WARNING, msg)
}

func (le *LogEntry) Crit(msg string) {
	le.Log(CRITICAL, msg)
}

func (le *LogEntry) Alert(msg string) {
	le.Log(ALERT, msg)
}

func (le *LogEntry) Emerg(msg string) {
	le.Log(EMERGENCY, msg)
}

func (le *LogEntry) Metric(msg string) {
	le.Payload["isMetric"] = true
	le.Log(EMERGENCY, msg)
}
