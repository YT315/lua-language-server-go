package capabililty

import (
	"context"
	"lualsp/protocol"
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

func (s *Server) DidSave(context.Context, *protocol.DidSaveTextDocumentParams) error {
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
