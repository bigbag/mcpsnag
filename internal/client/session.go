package client

import "github.com/bigbag/mcpsnag/internal/protocol"

type Session struct {
	ID           string
	Capabilities *protocol.ServerCapabilities
	ServerInfo   *protocol.Implementation
}

func (s *Session) IsValid() bool {
	return s.ID != ""
}
