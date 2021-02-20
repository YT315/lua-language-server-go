package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) DocumentLink(context.Context, *protocol.DocumentLinkParams) ([]protocol.DocumentLink /*DocumentLink[] | null*/, error) {
	return nil, nil
}

func (s *Server) ResolveDocumentLink(context.Context, *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, nil
}
