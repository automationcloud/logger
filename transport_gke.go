package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// TransportGKE implements Transport interface to work in k8s. It sends
// structured json payloads to stderr and stdout ready for consumption of
// fluentd logging agent running in a cluster.
type TransportGKE struct {
	logWriter   io.Writer
	errorWriter io.Writer
}

// NewTransportGKE creates a new Transport for GKE logging.
func NewTransportGKE() Transport {
	return &TransportGKE{os.Stdout, os.Stderr}
}

// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
type httpRequestGKE struct {
	RequestMethod string `json:"requestMethod,omitempty"`
	RequestURL    string `json:"requestUrl,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	RequestSize   string `json:"requestSize,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Referer       string `json:"referer,omitempty"`

	Status       int    `json:"status,omitempty"`
	ResponseSize string `json:"responseSize,omitempty"`
	RemoteIP     string `json:"remoteIp,omitempty"`
	ServerIP     string `json:"serverIp,omitempty"`
	Latency      string `json:"latency,omitempty"`

	CacheLookup    bool   `json:"cacheLookup,omitempty"`
	CacheHit       bool   `json:"cacheHit,omitempty"`
	CacheFillBytes string `json:"cacheFillBytes,omitempty"`
	// CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer,omitempty"`

}

// formatHTTPRequest expands HTTPRequest to be consumable by GKE logging agent.
func formatHTTPRequest(r *HTTPRequest) httpRequestGKE {
	return httpRequestGKE{
		RequestMethod:  r.Request.Method,
		RequestURL:     r.Request.URL.String(),
		Referer:        r.Request.Referer(),
		Protocol:       r.Request.Proto,
		RequestSize:    fmt.Sprintf("%d", r.RequestSize),
		Status:         r.Status,
		ResponseSize:   fmt.Sprintf("%d", r.ResponseSize),
		UserAgent:      r.Request.UserAgent(),
		RemoteIP:       r.RemoteIP,
		ServerIP:       r.ServerIP,
		Latency:        fmt.Sprintf("%fs", r.Latency.Seconds()),
		CacheLookup:    r.CacheLookup,
		CacheHit:       r.CacheHit,
		CacheFillBytes: r.CacheFillBytes,
		// CacheValidatedWithOriginServer: r.CacheValidatedWithOriginServer,
	}
}

// SendLog prepares log entry and sends it to stdout.
func (lt *TransportGKE) SendLog(le *LogEntry) error {
	payload := le.Payload
	if le.Message != "" {
		payload["message"] = le.Message
	}

	payload["severity"] = le.Severity
	if le.HTTPRequest != nil {
		payload["httpRequest"] = formatHTTPRequest(le.HTTPRequest)
	}

	if len(le.Labels) > 0 {
		payload["labels"] = le.Labels
	}

	return json.NewEncoder(lt.logWriter).Encode(payload)
}

// ReportError prepares error entry and sends it to stderr.
func (lt *TransportGKE) ReportError(le *ErrorEntry) error {
	return json.NewEncoder(lt.errorWriter).Encode(le)
}
