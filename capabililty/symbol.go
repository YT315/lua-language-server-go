package capabililty

import (
	"context"
	"lualsp/protocol"
)

//获取文章中的语法符号关系
func (s *Server) DocumentSymbol(context.Context, *protocol.DocumentSymbolParams) ([]interface{} /*SymbolInformation[] | DocumentSymbol[] | null*/, error) {
	return nil, nil
}
