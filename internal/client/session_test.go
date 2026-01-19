package client

import (
	"testing"

	"github.com/bigbag/mcpsnag/internal/protocol"
)

func TestSessionIsValid(t *testing.T) {
	tests := []struct {
		name     string
		session  Session
		expected bool
	}{
		{
			name:     "valid session with ID",
			session:  Session{ID: "test-session-123"},
			expected: true,
		},
		{
			name:     "invalid session empty ID",
			session:  Session{ID: ""},
			expected: false,
		},
		{
			name:     "valid session with capabilities",
			session:  Session{ID: "abc", Capabilities: &protocol.ServerCapabilities{}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSessionFields(t *testing.T) {
	caps := &protocol.ServerCapabilities{
		Tools: &protocol.ToolsCapability{ListChanged: true},
	}
	info := &protocol.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}

	session := Session{
		ID:           "session-xyz",
		Capabilities: caps,
		ServerInfo:   info,
	}

	if session.ID != "session-xyz" {
		t.Errorf("expected ID %q, got %q", "session-xyz", session.ID)
	}

	if session.Capabilities != caps {
		t.Error("Capabilities pointer mismatch")
	}

	if session.ServerInfo != info {
		t.Error("ServerInfo pointer mismatch")
	}

	if session.Capabilities.Tools == nil {
		t.Error("Capabilities.Tools should not be nil")
	}

	if !session.Capabilities.Tools.ListChanged {
		t.Error("Capabilities.Tools.ListChanged should be true")
	}

	if session.ServerInfo.Name != "test-server" {
		t.Errorf("expected ServerInfo.Name %q, got %q", "test-server", session.ServerInfo.Name)
	}
}
