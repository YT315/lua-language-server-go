package capabililty

import (
	"context"
	"lualsp/protocol"
)

//查找所有引用
func (s *Server) References(context.Context, *protocol.ReferenceParams) ([]protocol.Location /*Location[] | null*/, error) {
	return nil, nil
}
