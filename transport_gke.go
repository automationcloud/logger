package logger

import (
	"encoding/json"
	"fmt"
	"io"
)

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

type TransportGKE struct {
	logWriter   io.Writer
	errorWriter io.Writer
}

func (lt *TransportGKE) SendLog(le *LogEntry) {
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

	json.NewEncoder(lt.logWriter).Encode(payload)
}

func (lt *TransportGKE) ReportError(le *ErrorEntry) {
	json.NewEncoder(lt.errorWriter).Encode(le)
}
