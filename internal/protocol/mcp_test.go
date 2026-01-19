package protocol

import (
	"encoding/json"
	"testing"
)

func TestMCPVersion(t *testing.T) {
	if MCPVersion == "" {
		t.Error("MCPVersion should not be empty")
	}
}

func TestSessionHeader(t *testing.T) {
	expected := "Mcp-Session-Id"
	if SessionHeader != expected {
		t.Errorf("expected SessionHeader %q, got %q", expected, SessionHeader)
	}
}

func TestDefaultInitializeParams(t *testing.T) {
	params := DefaultInitializeParams()

	if params.ProtocolVersion != MCPVersion {
		t.Errorf("expected ProtocolVersion %q, got %q", MCPVersion, params.ProtocolVersion)
	}

	if params.ClientInfo.Name != "mcpsnag" {
		t.Errorf("expected ClientInfo.Name %q, got %q", "mcpsnag", params.ClientInfo.Name)
	}

	if params.ClientInfo.Version == "" {
		t.Error("ClientInfo.Version should not be empty")
	}

	if params.Capabilities.Roots == nil {
		t.Error("Capabilities.Roots should not be nil")
	}

	if !params.Capabilities.Roots.ListChanged {
		t.Error("Capabilities.Roots.ListChanged should be true")
	}
}

func TestInitializeParamsJSON(t *testing.T) {
	params := DefaultInitializeParams()

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed InitializeParams
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed.ProtocolVersion != params.ProtocolVersion {
		t.Errorf("ProtocolVersion mismatch: %q vs %q", parsed.ProtocolVersion, params.ProtocolVersion)
	}

	if parsed.ClientInfo.Name != params.ClientInfo.Name {
		t.Errorf("ClientInfo.Name mismatch: %q vs %q", parsed.ClientInfo.Name, params.ClientInfo.Name)
	}
}

func TestInitializeResultJSON(t *testing.T) {
	jsonData := `{
		"protocolVersion": "2025-03-26",
		"capabilities": {
			"tools": {"listChanged": true},
			"resources": {"subscribe": true, "listChanged": false}
		},
		"serverInfo": {
			"name": "test-server",
			"version": "1.0.0"
		}
	}`

	var result InitializeResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.ProtocolVersion != "2025-03-26" {
		t.Errorf("expected ProtocolVersion %q, got %q", "2025-03-26", result.ProtocolVersion)
	}

	if result.ServerInfo.Name != "test-server" {
		t.Errorf("expected ServerInfo.Name %q, got %q", "test-server", result.ServerInfo.Name)
	}

	if result.Capabilities.Tools == nil {
		t.Error("Capabilities.Tools should not be nil")
	}

	if !result.Capabilities.Tools.ListChanged {
		t.Error("Capabilities.Tools.ListChanged should be true")
	}

	if result.Capabilities.Resources == nil {
		t.Error("Capabilities.Resources should not be nil")
	}

	if !result.Capabilities.Resources.Subscribe {
		t.Error("Capabilities.Resources.Subscribe should be true")
	}
}

func TestServerCapabilitiesOptionalFields(t *testing.T) {
	jsonData := `{
		"protocolVersion": "2025-03-26",
		"capabilities": {},
		"serverInfo": {"name": "minimal", "version": "0.1"}
	}`

	var result InitializeResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Capabilities.Tools != nil {
		t.Error("Capabilities.Tools should be nil when not provided")
	}

	if result.Capabilities.Resources != nil {
		t.Error("Capabilities.Resources should be nil when not provided")
	}

	if result.Capabilities.Prompts != nil {
		t.Error("Capabilities.Prompts should be nil when not provided")
	}

	if result.Capabilities.Logging != nil {
		t.Error("Capabilities.Logging should be nil when not provided")
	}
}

func TestImplementationJSON(t *testing.T) {
	impl := Implementation{
		Name:    "test-client",
		Version: "2.0.0",
	}

	data, err := json.Marshal(impl)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed Implementation
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed.Name != impl.Name {
		t.Errorf("Name mismatch: %q vs %q", parsed.Name, impl.Name)
	}

	if parsed.Version != impl.Version {
		t.Errorf("Version mismatch: %q vs %q", parsed.Version, impl.Version)
	}
}
