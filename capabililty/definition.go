package capabililty

import (
	"context"
	"lualsp/protocol"
)

func (s *Server) TypeDefinition(context.Context, *protocol.TypeDefinitionParams) (protocol.Definition /*Definition | DefinitionLink[] | null*/, error) {
	return nil, nil
}
func (s *Server) Definition(context.Context, *protocol.DefinitionParams) (protocol.Definition /*Definition | DefinitionLink[] | null*/, error) {
	return nil, nil
}
