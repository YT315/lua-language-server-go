package capabililty

import (
	"context"
	"lualsp/protocol"
)

//not Implementation
func (s *Server) Implementation(context.Context, *protocol.ImplementationParams) (protocol.Definition /*Definition | DefinitionLink[] | null*/, error) {
	return nil, nil
}

func (s *Server) DocumentColor(context.Context, *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, nil
}

func (s *Server) ColorPresentation(context.Context, *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, nil
}

func (s *Server) Declaration(context.Context, *protocol.DeclarationParams) (protocol.Declaration /*Declaration | DeclarationLink[] | null*/, error) {
	return nil, nil
}

func (s *Server) SelectionRange(context.Context, *protocol.SelectionRangeParams) ([]protocol.SelectionRange /*SelectionRange[] | null*/, error) {
	return nil, nil
}

func (s *Server) PrepareCallHierarchy(context.Context, *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem /*CallHierarchyItem[] | null*/, error) {
	return nil, nil
}

func (s *Server) IncomingCalls(context.Context, *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall /*CallHierarchyIncomingCall[] | null*/, error) {
	return nil, nil
}

func (s *Server) OutgoingCalls(context.Context, *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall /*CallHierarchyOutgoingCall[] | null*/, error) {
	return nil, nil
}

func (s *Server) SignatureHelp(context.Context, *protocol.SignatureHelpParams) (*protocol.SignatureHelp /*SignatureHelp | null*/, error) {
	return nil, nil
}
func (s *Server) CodeAction(context.Context, *protocol.CodeActionParams) ([]protocol.CodeAction /*(Command | CodeAction)[] | null*/, error) {
	return nil, nil
}
func (s *Server) ResolveCodeAction(context.Context, *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, nil
}
func (s *Server) Symbol(context.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation /*SymbolInformation[] | null*/, error) {
	return nil, nil
}
func (s *Server) CodeLens(context.Context, *protocol.CodeLensParams) ([]protocol.CodeLens /*CodeLens[] | null*/, error) {
	return nil, nil
}
func (s *Server) ResolveCodeLens(context.Context, *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, nil
}
