package logger

// https://cloud.google.com/error-reporting/docs/formatting-error-messages
type ErrorEntryGKE struct {
	ServiceContext ServiceContext `json:"serviceContext"`
	Message        string         `json:"message"`
	Context        errorContext   `json:"context"`
	client         *Client
}

type errorContext struct {
	ReportLocation StackFrame         `json:"reportLocation"`
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
