package capabililty

import (
	"context"
	"lualsp/protocol"
)

//高亮此位置符号所有的引用以及定义
func (s *Server) DocumentHighlight(context.Context, *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight /*DocumentHighlight[] | null*/, error) {
	return nil, nil
}
