package protocol

import (
	"context"
	"errors"
	"lualsp/auxiliary"

	"github.com/sourcegraph/jsonrpc2"
)

//NewHandler 事件
func NewHandler(server Server) jsonrpc2.Handler {
	return jsonrpc2.HandlerWithError(protocolHandler{
		server: server,
	}.handle)
}

type protocolHandler struct {
	server Server
}

func (h protocolHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	//将客户端传递给下文
	ctx = context.WithValue(ctx, auxiliary.CtxKey("cllent"), clientDispatcher{Conn: conn})
	ok, err := serverDispatch(ctx, h.server,
		func(ctx context.Context, data interface{}, err error) error {
			result = data
			return nil
		},
		*req)
	if !ok {
		err = errors.New("不支持此函数:" + req.Method)
	}
	return
}
