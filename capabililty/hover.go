package capabililty

import (
	"context"
	"lualsp/protocol"
)

//浮窗显示信息
func (s *Server) Hover(context.Context, *protocol.HoverParams) (*protocol.Hover /*Hover | null*/, error) {
	return nil, nil
}
