package logger

import (
	"net/http"
	"testing"
)

func TestWithPayload(t *testing.T) {
	e := client.NewLogEntry().WithPayload("hello", "world")
	if e.Payload["hello"] != "world" {
		t.Error("expect payload to keep dictionary of arbitrary values")
	}
}

func TestWithRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	e := client.NewLogEntry().WithRequest(&HTTPRequest{
		Request: r,
	})

	if e.HTTPRequest.Request.URL.Path != "/" {
		t.Error("expect payload to keep information about request")
	}
}

func TestLogMethods(t *testing.T) {
	cases := map[Severity]func(){
		INFO:      func() { client.NewLogEntry().Info("message text") },
		DEBUG:     func() { client.NewLogEntry().Debug("message text") },
		CRITICAL:  func() { client.NewLogEntry().Crit("message text") },
		ALERT:     func() { client.NewLogEntry().Alert("message text") },
		WARNING:   func() { client.NewLogEntry().Warn("message text") },
		EMERGENCY: func() { client.NewLogEntry().Emerg("message text") },
	}

	for expectedSeverity, fn := range cases {
		fn()
		actualSeverity := DefaultTransport.(*mockTransport).logEntry.Severity
		if actualSeverity != expectedSeverity {
			t.Errorf(
				"expected %v severity, got %v",
				expectedSeverity,
				actualSeverity,
			)
		}
	}
}

func TestLogMetric(t *testing.T) {
	client.NewLogEntry().Metric("hello")
	if DefaultTransport.(*mockTransport).logEntry.Payload["isMetric"] != true {
		t.Error("expect payload to have isMetric flag")
	}
}
