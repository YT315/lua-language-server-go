package capabililty

import (
	"context"
	"lualsp/auxiliary"
	"lualsp/protocol"

	"lualsp/logger"
)

//文本变化
func (s *Server) DidOpen(context.Context, *protocol.DidOpenTextDocumentParams) error {
	return nil
}

func (s *Server) DidChange(context.Context, *protocol.DidChangeTextDocumentParams) error {
	return nil
}

func (s *Server) DidClose(context.Context, *protocol.DidCloseTextDocumentParams) error {
	return nil
}

func (s *Server) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	data := ctx.Value(auxiliary.CtxKey("client"))
	client, ok := data.(protocol.Client)
	if ok {
		pam := protocol.PublishDiagnosticsParams{URI: params.TextDocument.URI}
		diag := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line: 1, Character: 0,
				},
				End: protocol.Position{
					Line: 1, Character: 10,
				},
			},
			Severity: protocol.SeverityError,
			Code:     "test_Code",
			Source:   "test_Source",
			Message:  "test_Message",
			Tags:     []protocol.DiagnosticTag{protocol.Deprecated},
			Data:     "test_Data",
			RelatedInformation: []protocol.DiagnosticRelatedInformation{
				protocol.DiagnosticRelatedInformation{
					Location: protocol.Location{
						URI: params.TextDocument.URI,
						Range: protocol.Range{
							Start: protocol.Position{
								Line: 4, Character: 0,
							},
							End: protocol.Position{
								Line: 4, Character: 10,
							},
						},
					},
					Message: "test_DiagnosticRelatedInformation",
				},
			},
			CodeDescription: &protocol.CodeDescription{
				Href: "www.baidu.com",
			},
		}
		pam.Diagnostics = append(pam.Diagnostics, diag)
		//client.WorkspaceFolders(ctx)
		//logger.Debugln(res)
		client.PublishDiagnostics(ctx, &pam)

	} else {
		logger.Debugln("fail")
	}

	return nil
}

func (s *Server) WillSave(context.Context, *protocol.WillSaveTextDocumentParams) error {
	return nil
}
func (s *Server) WillSaveWaitUntil(context.Context, *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit /*TextEdit[] | null*/, error) {
	return nil, nil
}

func (s *Server) DidChangeWorkspaceFolders(context.Context, *protocol.DidChangeWorkspaceFoldersParams) error {
	return nil
}

func (s *Server) DidChangeWatchedFiles(context.Context, *protocol.DidChangeWatchedFilesParams) error {
	return nil
}
