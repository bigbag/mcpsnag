package protocol

import "encoding/json"

const JSONRPCVersion = "2.0"

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func NewRequest(id any, method string, params any) (*Request, error) {
	var paramsRaw json.RawMessage
	if params != nil {
		p, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsRaw = p
	}
	return &Request{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  method,
		Params:  paramsRaw,
	}, nil
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func NewNotification(method string, params any) (*Notification, error) {
	var paramsRaw json.RawMessage
	if params != nil {
		p, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsRaw = p
	}
	return &Notification{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  paramsRaw,
	}, nil
}

type UserRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}
