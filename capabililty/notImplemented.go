package capabililty

import (
	"context"
	"lualsp/protocol"
)

//转到实现处
func (s *Server) Implementation(context.Context, *protocol.ImplementationParams) (protocol.Definition /*Definition | DefinitionLink[] | null*/, error) {
	return nil, nil
}

//出现颜色选项
func (s *Server) DocumentColor(context.Context, *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, nil
}

func (s *Server) ColorPresentation(context.Context, *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, nil
}

//转到声明处
func (s *Server) Declaration(context.Context, *protocol.DeclarationParams) (protocol.Declaration /*Declaration | DeclarationLink[] | null*/, error) {
	return nil, nil
}

//推荐选择范围
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

//界面显示小灯泡,选择命令
func (s *Server) CodeAction(ctx context.Context, param *protocol.CodeActionParams) ([]protocol.CodeAction /*(Command | CodeAction)[] | null*/, error) {
	if len(param.Context.Diagnostics) != 0 {

		tx := protocol.TextEdit{
			Range: protocol.Range{
				Start: protocol.Position{
					Line: 3, Character: 0,
				},
				End: protocol.Position{
					Line: 3, Character: 10,
				},
			},
			NewText: "hello my baby",
		}
		ca := protocol.CodeAction{
			Title:       "test_codeaction",
			Kind:        protocol.QuickFix,
			Diagnostics: param.Context.Diagnostics,

			Edit: &protocol.WorkspaceEdit{
				Changes: make(map[string][]protocol.TextEdit),
			},
		}
		ca.Edit.Changes[string(param.TextDocument.URI)] = []protocol.TextEdit{tx}

		return []protocol.CodeAction{ca}, nil
	}
	return nil, nil
}
func (s *Server) ResolveCodeAction(context.Context, *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, nil
}

func (s *Server) Symbol(context.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation /*SymbolInformation[] | null*/, error) {
	return nil, nil
}

//界面显示文字,执行命令
func (s *Server) CodeLens(context.Context, *protocol.CodeLensParams) ([]protocol.CodeLens /*CodeLens[] | null*/, error) {
	return nil, nil
}
func (s *Server) ResolveCodeLens(context.Context, *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, nil
}

func (s *Server) Moniker(context.Context, *protocol.MonikerParams) ([]protocol.Moniker /*Moniker[] | null*/, error) {
	return nil, nil
}
func (s *Server) NonstandardRequest(ctx context.Context, method string, params interface{}) (interface{}, error) {
	return nil, nil
}
func (s *Server) SetTrace(context.Context, *protocol.SetTraceParams) error {
	return nil
}
func (s *Server) LogTrace(context.Context, *protocol.LogTraceParams) error {
	return nil
}
