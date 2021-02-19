package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) DocumentSymbol(context.Context, *protocol.DocumentSymbolParams) ([]interface{} /*SymbolInformation[] | DocumentSymbol[] | null*/, error) {
	return nil, nil
}
