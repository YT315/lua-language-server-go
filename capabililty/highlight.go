package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) DocumentHighlight(context.Context, *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight /*DocumentHighlight[] | null*/, error) {
	return nil, nil
}
