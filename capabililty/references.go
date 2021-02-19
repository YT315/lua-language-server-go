package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) References(context.Context, *protocol.ReferenceParams) ([]protocol.Location /*Location[] | null*/, error) {
	return nil, nil
}
