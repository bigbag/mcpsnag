package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bigbag/mcpsnag/internal/protocol"
)

type Transport struct {
	endpoint   string
	httpClient *http.Client
	headers    map[string]string
}

func NewTransport(endpoint string, timeout time.Duration) *Transport {
	return &Transport{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: make(map[string]string),
	}
}

func (t *Transport) SetHeader(key, value string) {
	t.headers[key] = value
}

func (t *Transport) Post(body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", t.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	return t.httpClient.Do(req)
}

type SSEEvent struct {
	Event string
	Data  string
	ID    string
}

func ParseSSEStream(r io.Reader, handler func(event SSEEvent) error) error {
	scanner := bufio.NewScanner(r)
	var current SSEEvent

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if current.Data != "" {
				if err := handler(current); err != nil {
					return err
				}
			}
			current = SSEEvent{}
			continue
		}

		if strings.HasPrefix(line, ":") {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		field := line[:colonIdx]
		value := strings.TrimPrefix(line[colonIdx+1:], " ")

		switch field {
		case "event":
			current.Event = value
		case "data":
			if current.Data != "" {
				current.Data += "\n"
			}
			current.Data += value
		case "id":
			current.ID = value
		}
	}

	if current.Data != "" {
		if err := handler(current); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (t *Transport) PostAndReadResponse(body []byte, stream bool, onEvent func(protocol.Response) error) (*protocol.Response, string, error) {
	resp, err := t.Post(body)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	sessionID := resp.Header.Get(protocol.SessionHeader)

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, sessionID, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil, sessionID, nil
	}

	if strings.HasPrefix(contentType, "text/event-stream") {
		var lastResponse *protocol.Response
		err := ParseSSEStream(resp.Body, func(event SSEEvent) error {
			if event.Event == "message" || event.Event == "" {
				var jsonResp protocol.Response
				if err := json.Unmarshal([]byte(event.Data), &jsonResp); err != nil {
					return err
				}
				lastResponse = &jsonResp
				if stream && onEvent != nil {
					return onEvent(jsonResp)
				}
			}
			return nil
		})
		if err != nil {
			return nil, sessionID, err
		}
		return lastResponse, sessionID, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, sessionID, err
	}

	var jsonResp protocol.Response
	if err := json.Unmarshal(bodyBytes, &jsonResp); err != nil {
		return nil, sessionID, fmt.Errorf("invalid JSON response: %w", err)
	}

	return &jsonResp, sessionID, nil
}
