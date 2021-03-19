package capabililty

import (
	"context"
	"lualsp/protocol"
)

//格式化
func (s *Server) Formatting(context.Context, *protocol.DocumentFormattingParams) ([]protocol.TextEdit /*TextEdit[] | null*/, error) {
	return nil, nil
}
func (s *Server) RangeFormatting(context.Context, *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit /*TextEdit[] | null*/, error) {
	return nil, nil
}
func (s *Server) OnTypeFormatting(context.Context, *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit /*TextEdit[] | null*/, error) {
	return nil, nil
}
