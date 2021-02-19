package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) SemanticTokensFull(context.Context, *protocol.SemanticTokensParams) (*protocol.SemanticTokens /*SemanticTokens | null*/, error) {
	return nil, nil
}
func (s *Server) SemanticTokensFullDelta(context.Context, *protocol.SemanticTokensDeltaParams) (interface{} /* SemanticTokens | SemanticTokensDelta | nil*/, error) {
	return nil, nil
}
func (s *Server) SemanticTokensRange(context.Context, *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens /*SemanticTokens | null*/, error) {
	return nil, nil
}
func (s *Server) SemanticTokensRefresh(context.Context) error {
	return nil
}
