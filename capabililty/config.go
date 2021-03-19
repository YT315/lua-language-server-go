package capabililty

import (
	"context"
	"lualsp/protocol"
)

//配置修改
func (s *Server) DidChangeConfiguration(context.Context, *protocol.DidChangeConfigurationParams) error {
	return nil
}
