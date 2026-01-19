package client

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bigbag/mcpsnag/internal/protocol"
)

type Client struct {
	transport *Transport
	session   *Session
	requestID atomic.Int64
	stream    bool
}

type Options struct {
	Endpoint  string
	Headers   map[string]string
	SessionID string
	Timeout   time.Duration
	Stream    bool
}

func New(opts Options) *Client {
	t := NewTransport(opts.Endpoint, opts.Timeout)
	for k, v := range opts.Headers {
		t.SetHeader(k, v)
	}

	session := &Session{ID: opts.SessionID}
	if session.ID != "" {
		t.SetHeader(protocol.SessionHeader, session.ID)
	}

	return &Client{
		transport: t,
		session:   session,
		stream:    opts.Stream,
	}
}

func (c *Client) nextID() int64 {
	return c.requestID.Add(1)
}

func (c *Client) Initialize() (*protocol.InitializeResult, error) {
	params := protocol.DefaultInitializeParams()
	req, err := protocol.NewRequest(c.nextID(), "initialize", params)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, sessionID, err := c.transport.PostAndReadResponse(body, false, nil)
	if err != nil {
		return nil, err
	}

	if sessionID != "" {
		c.session.ID = sessionID
		c.transport.SetHeader(protocol.SessionHeader, sessionID)
	}

	if resp == nil {
		return nil, fmt.Errorf("no response from initialize")
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	var result protocol.InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse initialize result: %w", err)
	}

	c.session.Capabilities = &result.Capabilities
	c.session.ServerInfo = &result.ServerInfo

	notif, err := protocol.NewNotification("notifications/initialized", nil)
	if err != nil {
		return nil, err
	}

	notifBody, err := json.Marshal(notif)
	if err != nil {
		return nil, err
	}

	_, _, err = c.transport.PostAndReadResponse(notifBody, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send initialized notification: %w", err)
	}

	return &result, nil
}

func (c *Client) Session() *Session {
	return c.session
}

func (c *Client) Request(method string, params json.RawMessage, onEvent func(protocol.Response) error) (*protocol.Response, error) {
	req, err := protocol.NewRequest(c.nextID(), method, nil)
	if err != nil {
		return nil, err
	}
	req.Params = params

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.transport.PostAndReadResponse(body, c.stream, onEvent)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Error != nil {
		return resp, resp.Error
	}

	return resp, nil
}

func (c *Client) RawRequest(body []byte, onEvent func(protocol.Response) error) (*protocol.Response, string, error) {
	return c.transport.PostAndReadResponse(body, c.stream, onEvent)
}
