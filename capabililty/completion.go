package capabililty

import (
	"context"
	"lualsp/auxiliary"
	"lualsp/logger"
	"lualsp/protocol"
)

//补全提示
func (s *Server) Completion(ctx context.Context, param *protocol.CompletionParams) (*protocol.CompletionList /*CompletionItem[] | CompletionList | null*/, error) {
	data := ctx.Value(auxiliary.CtxKey("client"))
	client, ok := data.(protocol.Client)
	if ok {
		client.WorkspaceFolders(ctx)
		//logger.Debugln(res)
	} else {
		logger.Debugln("fail")
	}

	return nil, nil
}
func (s *Server) Resolve(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, nil
}
