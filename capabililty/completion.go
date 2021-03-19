package capabililty

import (
	"context"
	"lualsp/protocol"
)

//补全提示
func (s *Server) Completion(context.Context, *protocol.CompletionParams) (*protocol.CompletionList /*CompletionItem[] | CompletionList | null*/, error) {
	return nil, nil
}
func (s *Server) Resolve(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, nil
}
