package main

import (
	"context"
	"flag"
	"fmt"
	"lualsp/auxiliary"
	"lualsp/capabililty"
	"lualsp/logger"
	"lualsp/protocol"
	"lualsp/syntax"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

const version = "0.1"

var (
	ast      = flag.String("AST", "", "输出文件的抽象语法树并退出")
	mode     = flag.String("mode", "stdio", "与客户端的通信模式 stdio|tcp|websocket")
	addr     = flag.String("addr", ":7529", "服务器监听地址在stdio和websocket模式下")
	logLevel = flag.String("logLevel", "none", "设置日志等级 debug|info|warning|error|none")
	logWay   = flag.String("logWay", "file", "设置日志输出方式,当mode是stdio时,此项不可以设置stdio, file|stdio|all")
	vsrsion  = flag.Bool("version", false, "输出版本号并退出")
)

func main() {
	flag.Parse()
	logger.Init(*logLevel, *logWay)
	if *vsrsion {
		fmt.Println(version)
		return
	}
	if *ast != "" {
		f, err := os.Open(*ast)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer func() {
			if err = f.Close(); err != nil {
				fmt.Println(err.Error())
			}
		}()

		lex := syntax.NewLexer(f, func(line, col uint, msg string) {
			println("err:- line:", line, "col:", col, "msg:", msg)
		})
		lex.Parse()
		auxiliary.DrawTree(lex.Block, "All")
		return
	}

	switch *mode {
	case "stdio":
		logger.Debugln("star server in stdio mode")
		server := capabililty.NewServer()
		handler := protocol.NewHandler(server)
		<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}), handler).DisconnectNotify()
		logger.Debugln("connection closed")
	default:
		logger.Errorln("invalid mode" + *mode)
	}
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
