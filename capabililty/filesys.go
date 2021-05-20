package capabililty

import (
	"context"
	"lualsp/auxiliary"
	"lualsp/logger"
	"lualsp/protocol"
)

//文本变化
func (s *Server) DidOpen(context.Context, *protocol.DidOpenTextDocumentParams) error {
	return nil
}

func (s *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	path := auxiliary.UriToPath(string(params.TextDocument.URI))
	exist := false
	if len(params.ContentChanges) > 1 {
		logger.Debugf("ddd")
	}
	for _, w := range s.project.Workspaces {
		if file, ok := w.Files[path]; ok {
			for _, change := range params.ContentChanges {
				logger.Debugln(change.Text)
				file.Content.Insert(
					int(change.Range.Start.Line),
					int(change.Range.Start.Character),
					int(change.Range.End.Line),
					int(change.Range.End.Character),
					int(change.RangeLength),
					change.Text,
				)
			}
			exist = true
			break
		}
	}
	if !exist {
		return nil
	}
	return nil
}

func (s *Server) DidClose(context.Context, *protocol.DidCloseTextDocumentParams) error {
	//不需要反应
	return nil
}

func (s *Server) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	//暂时用于测试事件
	data := ctx.Value(auxiliary.CtxKey("client"))
	client, ok := data.(protocol.Client)
	if ok {
		for _, ws := range s.project.Workspaces {
			for uri, file := range ws.Files {
				if len(file.Diagnostics) > 0 {
					client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
						URI:         protocol.DocumentURI(uri),
						Diagnostics: []protocol.Diagnostic{},
					})
				}

			}
		}
	}
	// data := ctx.Value(auxiliary.CtxKey("client"))
	// client, ok := data.(protocol.Client)
	// if ok {
	// 	pam := protocol.PublishDiagnosticsParams{URI: params.TextDocument.URI}
	// 	diag := protocol.Diagnostic{
	// 		Range: protocol.Range{
	// 			Start: protocol.Position{
	// 				Line: 1, Character: 0,
	// 			},
	// 			End: protocol.Position{
	// 				Line: 1, Character: 10,
	// 			},
	// 		},
	// 		Severity: protocol.SeverityError,
	// 		Code:     "test_Code",
	// 		Source:   "test_Source",
	// 		Message:  "test_Message",
	// 		Tags:     []protocol.DiagnosticTag{protocol.Deprecated},
	// 		Data:     "test_Data",
	// 		RelatedInformation: []protocol.DiagnosticRelatedInformation{
	// 			protocol.DiagnosticRelatedInformation{
	// 				Location: protocol.Location{
	// 					URI: params.TextDocument.URI,
	// 					Range: protocol.Range{
	// 						Start: protocol.Position{
	// 							Line: 4, Character: 0,
	// 						},
	// 						End: protocol.Position{
	// 							Line: 4, Character: 10,
	// 						},
	// 					},
	// 				},
	// 				Message: "test_DiagnosticRelatedInformation",
	// 			},
	// 		},
	// 		CodeDescription: &protocol.CodeDescription{
	// 			Href: "www.baidu.com",
	// 		},
	// 	}
	// 	pam.Diagnostics = append(pam.Diagnostics, diag)
	// 	//client.WorkspaceFolders(ctx)
	// 	//logger.Debugln(res)
	// 	client.PublishDiagnostics(ctx, &pam)

	// } else {
	// 	logger.Debugln("fail")
	// }

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
