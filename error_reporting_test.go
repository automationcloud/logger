package logger

type FakeWriter struct {
	Data []byte
}

func (w *FakeWriter) Write(p []byte) (n int, err error) {
	w.Data = p
	return 0, nil
}

/*
func TestReportError(t *testing.T) {
	writer := &FakeWriter{}

	reporter := NewLogger("workers-autoscaler", "test-version")
	reporter.transport

	func() {
		reporter.LogError(errors.New("test error"))
	}()

	var structuredError StructuredError
	err := json.Unmarshal(writer.Data, &structuredError)
	if err != nil {
		t.FailNow()
	}

	equal(t, structuredError.ServiceContext.Service, "workers-autoscaler")
	equal(t, structuredError.ServiceContext.Version, "test-version")
	equal(t, structuredError.Context.ReportLocation.FunctionName, "TestReportError.func1")
	if !strings.HasSuffix(structuredError.Context.ReportLocation.FilePath, "error_reporting_test.go") {
		t.Errorf(
			"Expected reportLocation.filePath to end with %v, got %v",
			"error_reporting_test.go",
			structuredError.Context.ReportLocation.FilePath,
		)
	}

	if !strings.HasPrefix(structuredError.Message, "test error\n") {
		t.Errorf(
			"Expected message to end with %v, got %v",
			"test error\n",
			structuredError.Message,
		)
	}
}

func equal(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("Expected %v, got %v", b, a)
	}
}
*/
