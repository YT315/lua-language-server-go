package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) Hover(context.Context, *protocol.HoverParams) (*protocol.Hover /*Hover | null*/, error) {
	return nil, nil
}
