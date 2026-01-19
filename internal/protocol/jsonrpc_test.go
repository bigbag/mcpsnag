package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, err := NewRequest(1, "test/method", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	if req.JSONRPC != JSONRPCVersion {
		t.Errorf("expected JSONRPC %q, got %q", JSONRPCVersion, req.JSONRPC)
	}
	if req.ID != 1 {
		t.Errorf("expected ID 1, got %v", req.ID)
	}
	if req.Method != "test/method" {
		t.Errorf("expected method %q, got %q", "test/method", req.Method)
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed Request
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed.Method != req.Method {
		t.Errorf("roundtrip method mismatch: %q vs %q", parsed.Method, req.Method)
	}
}

func TestNewRequestNilParams(t *testing.T) {
	req, err := NewRequest(1, "test/method", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	if req.Params != nil {
		t.Errorf("expected nil params, got %v", req.Params)
	}
}

func TestNewNotification(t *testing.T) {
	notif, err := NewNotification("notifications/test", nil)
	if err != nil {
		t.Fatalf("NewNotification failed: %v", err)
	}

	if notif.JSONRPC != JSONRPCVersion {
		t.Errorf("expected JSONRPC %q, got %q", JSONRPCVersion, notif.JSONRPC)
	}
	if notif.Method != "notifications/test" {
		t.Errorf("expected method %q, got %q", "notifications/test", notif.Method)
	}
}

func TestResponseWithError(t *testing.T) {
	jsonData := `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`

	var resp Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != -32600 {
		t.Errorf("expected code -32600, got %d", resp.Error.Code)
	}
	if resp.Error.Message != "Invalid Request" {
		t.Errorf("expected message %q, got %q", "Invalid Request", resp.Error.Message)
	}
}

func TestErrorImplementsError(t *testing.T) {
	e := &Error{Code: -32600, Message: "Invalid Request"}
	var _ error = e

	if e.Error() != "Invalid Request" {
		t.Errorf("expected %q, got %q", "Invalid Request", e.Error())
	}
}
