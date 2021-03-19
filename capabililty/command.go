package capabililty

import (
	"context"
	"lualsp/protocol"
)

//服务器初始化时,会通过(ExecuteCommandProvider)告诉客户端,自己有哪些命令,客户端执行命令,就会回调此函数
func (s *Server) ExecuteCommand(context.Context, *protocol.ExecuteCommandParams) (interface{} /*any | null*/, error) {
	return nil, nil
}
