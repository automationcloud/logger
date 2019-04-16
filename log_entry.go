package logger

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
	Payload     map[string]interface{}
	logger      *Logger
}

// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
type HTTPRequest struct {
	// The request method. Examples: "GET", "HEAD", "PUT", "POST".
	RequestMethod string `json:"requestMethod,omitempty"`

	// The scheme (http, https), the host name, the path and the query portion
	// of the URL that was requested. Example:
	// "http://example.com/some/info?color=red".
	RequestURL string `json:"requestUrl,omitempty"`

	// The size of the HTTP request message in bytes, including the request headers and the request body.
	RequestSize string `json:"requestSize,omitempty"`

	// The response code indicating the status of response. Examples: 200, 404.
	Status int `json:"status,omitempty"`

	// The size of the HTTP response message sent back to the client, in bytes, including the response headers and the response body.
	ResponseSize string `json:"responseSize,omitempty"`
	UserAgent    string `json:"userAgent,omitempty"`
	RemoteIP     string `json:"remoteIp,omitempty"`
	ServerIP     string `json:"serverIp,omitempty"`
	Referer      string `json:"referer,omitempty"`
	Latency      string `json:"latency,omitempty"`

	CacheLookup                    bool   `json:"cacheLookup,omitempty"`
	CacheHit                       bool   `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 string `json:"cacheFillBytes,omitempty"`

	Protocol string `json:"protocol,omitempty"`
}

func (le *LogEntry) WithPayload(key string, val interface{}) *LogEntry {
	le.Payload[key] = val
	return le
}

func (le *LogEntry) WithPayloadHTTP(hr *HTTPRequest) *LogEntry {
	le.HTTPRequest = hr
	return le
}

func (le *LogEntry) Log(s Severity, msg string) {
	le.Message = msg
	le.Severity = s
	le.logger.transport.SendLog(le)
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
