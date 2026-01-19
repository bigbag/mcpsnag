package client

import (
	"strings"
	"testing"
)

func TestParseSSEStream(t *testing.T) {
	input := `event: message
data: {"jsonrpc":"2.0","id":1,"result":{"tools":[]}}

event: message
data: {"jsonrpc":"2.0","id":2,"result":{"done":true}}

`
	var events []SSEEvent
	err := ParseSSEStream(strings.NewReader(input), func(event SSEEvent) error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		t.Fatalf("ParseSSEStream failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].Event != "message" {
		t.Errorf("expected event 'message', got %q", events[0].Event)
	}
	if !strings.Contains(events[0].Data, `"tools":[]`) {
		t.Errorf("unexpected data: %s", events[0].Data)
	}
}

func TestParseSSEStreamMultilineData(t *testing.T) {
	input := `event: message
data: line1
data: line2

`
	var events []SSEEvent
	err := ParseSSEStream(strings.NewReader(input), func(event SSEEvent) error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		t.Fatalf("ParseSSEStream failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	expected := "line1\nline2"
	if events[0].Data != expected {
		t.Errorf("expected data %q, got %q", expected, events[0].Data)
	}
}

func TestParseSSEStreamWithComments(t *testing.T) {
	input := `: this is a comment
event: message
data: test

`
	var events []SSEEvent
	err := ParseSSEStream(strings.NewReader(input), func(event SSEEvent) error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		t.Fatalf("ParseSSEStream failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Data != "test" {
		t.Errorf("expected data 'test', got %q", events[0].Data)
	}
}

func TestParseSSEStreamWithID(t *testing.T) {
	input := `event: message
id: 123
data: test

`
	var events []SSEEvent
	err := ParseSSEStream(strings.NewReader(input), func(event SSEEvent) error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		t.Fatalf("ParseSSEStream failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].ID != "123" {
		t.Errorf("expected ID '123', got %q", events[0].ID)
	}
}

func TestParseSSEStreamNoEventType(t *testing.T) {
	input := `data: {"result":"ok"}

`
	var events []SSEEvent
	err := ParseSSEStream(strings.NewReader(input), func(event SSEEvent) error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		t.Fatalf("ParseSSEStream failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Event != "" {
		t.Errorf("expected empty event type, got %q", events[0].Event)
	}
}
