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
		res, _ := client.WorkspaceFolders(ctx)
		logger.Debugln(res)
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
