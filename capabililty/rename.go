package capabililty

import (
	"context"
	"lualsp/protocol"
)

//修改名字
func (s *Server) Rename(context.Context, *protocol.RenameParams) (*protocol.WorkspaceEdit /*WorkspaceEdit | null*/, error) {
	return nil, nil
}

func (s *Server) PrepareRename(context.Context, *protocol.PrepareRenameParams) (*protocol.Range /*Range | { range: Range, placeholder: string } | { defaultBehavior: boolean } | null*/, error) {
	return nil, nil
}
