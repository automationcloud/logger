package logger

import "testing"

func TestWithPayload(t *testing.T) {
	e := logger.NewLogEntry().WithPayload("hello", "world")
	if e.payload["hello"] != "world" {
		t.Error("expect payload to keep dictionary of arbitrary values")
	}
}

func TestWithPayloadHTTP(t *testing.T) {
	e := logger.NewLogEntry().WithPayloadHTTP(&HTTPRequest{
		RequestURL:    "/",
		RequestMethod: "GET",
	})

	if e.HTTPRequest.RequestURL != "/" {
		t.Error("expect payload to keep information about request")
	}
}

func TestLogMethods(t *testing.T) {
	cases := map[Severity]func(){
		INFO:      func() { logger.NewLogEntry().Info("message text") },
		DEBUG:     func() { logger.NewLogEntry().Debug("message text") },
		CRITICAL:  func() { logger.NewLogEntry().Crit("message text") },
		ALERT:     func() { logger.NewLogEntry().Alert("message text") },
		WARNING:   func() { logger.NewLogEntry().Warn("message text") },
		EMERGENCY: func() { logger.NewLogEntry().Emerg("message text") },
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
	logger.NewLogEntry().Metric("hello")
	if DefaultTransport.(*mockTransport).logEntry.payload["isMetric"] != true {
		t.Error("expect payload to have isMetric flag")
	}
}
