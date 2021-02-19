package capabililty

import (
	"context"
	"lualsp/protocol"
)

//文件折叠
func (s *Server) FoldingRange(context.Context, *protocol.FoldingRangeParams) ([]protocol.FoldingRange /*FoldingRange[] | null*/, error) {
	return nil, nil
}
