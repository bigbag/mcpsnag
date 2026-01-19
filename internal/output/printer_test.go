package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrinterPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	data := map[string]string{"key": "value"}
	err := p.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\"key\"") {
		t.Errorf("expected key in output, got %s", output)
	}
	if !strings.Contains(output, "\"value\"") {
		t.Errorf("expected value in output, got %s", output)
	}
}

func TestPrinterPrintJSONCompact(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, true, false)

	data := map[string]string{"key": "value"}
	err := p.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	expected := `{"key":"value"}`
	if output != expected {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestPrinterPrintJSONPretty(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	data := map[string]string{"key": "value"}
	err := p.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\n") {
		t.Errorf("expected pretty output with newlines, got %s", output)
	}
}

func TestPrinterPrintRawJSON(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	raw := json.RawMessage(`{"tools":[{"name":"test"}]}`)
	err := p.PrintRawJSON(raw)
	if err != nil {
		t.Fatalf("PrintRawJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "tools") {
		t.Errorf("expected tools in output, got %s", output)
	}
}

func TestPrinterPrintRawJSONCompact(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, true, false)

	raw := json.RawMessage(`{"key":"value"}`)
	err := p.PrintRawJSON(raw)
	if err != nil {
		t.Fatalf("PrintRawJSON failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	expected := `{"key":"value"}`
	if output != expected {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestPrinterPrintRawJSONInvalid(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	raw := json.RawMessage(`not valid json`)
	err := p.PrintRawJSON(raw)
	if err != nil {
		t.Fatalf("PrintRawJSON should not fail on invalid JSON: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "not valid json" {
		t.Errorf("expected raw output for invalid JSON, got %s", output)
	}
}

func TestPrinterPrintVerbose(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, true)

	p.PrintVerbose("test message %s", "arg")

	output := buf.String()
	if !strings.Contains(output, "test message arg") {
		t.Errorf("expected verbose message, got %s", output)
	}
}

func TestPrinterPrintVerboseDisabled(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	p.PrintVerbose("test message")

	output := buf.String()
	if output != "" {
		t.Errorf("expected no output when verbose disabled, got %s", output)
	}
}

func TestPrinterPrintError(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	p.PrintError(nil)
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error prefix, got %s", output)
	}
}

func TestPrinterPrintSessionInfo(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	p.PrintSessionInfo("test-session-123")

	output := buf.String()
	if !strings.Contains(output, "sessionId") {
		t.Errorf("expected sessionId in output, got %s", output)
	}
	if !strings.Contains(output, "test-session-123") {
		t.Errorf("expected session value in output, got %s", output)
	}
}

func TestPrinterPrintRequestVerbose(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, true)

	headers := map[string]string{"Authorization": "Bearer token"}
	body := []byte(`{"method":"test"}`)

	p.PrintRequest("POST", "http://localhost/mcp", headers, body)

	output := buf.String()
	if !strings.Contains(output, "POST") {
		t.Errorf("expected POST in output, got %s", output)
	}
	if !strings.Contains(output, "Authorization") {
		t.Errorf("expected header in output, got %s", output)
	}
}

func TestPrinterPrintRequestNotVerbose(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, false, false)

	p.PrintRequest("POST", "http://localhost/mcp", nil, nil)

	output := buf.String()
	if output != "" {
		t.Errorf("expected no output when not verbose, got %s", output)
	}
}
