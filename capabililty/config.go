package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) DidChangeConfiguration(context.Context, *protocol.DidChangeConfigurationParams) error {
	return nil
}
