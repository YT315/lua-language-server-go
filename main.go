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
	"net"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

const version = "0.1"

var (
	ast      = flag.String("AST", "", "输出文件的抽象语法树并退出")
	mode     = flag.String("mode", "stdio", "与客户端的通信模式 stdio|tcp")
	addr     = flag.String("addr", ":61358", "服务器监听地址在stdio和websocket模式下")
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

	case "tcp":
		logger.Debugln("star server in tcp mode")
		listener, err := net.Listen("tcp", *addr)
		if err != nil {
			logger.Errorln(err.Error())
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Errorln(err.Error())
			}
			logger.Debugln("received incoming connection ", conn.RemoteAddr())
			server := capabililty.NewServer()
			handler := protocol.NewHandler(server)
			jsonrpc2Connection := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), handler)
			go func() {
				<-jsonrpc2Connection.DisconnectNotify()
				defer jsonrpc2Connection.Close()
				logger.Debugln("connection", conn.RemoteAddr(), "closed")
			}()
		}

	default:
		logger.Errorln("invalid mode" + *mode)
	}
}

/*
func logTrace() jsonrpc2.ConnOpt {
	return func(c *jsonrpc2.Conn) {
		// Remember reqs we have received so we can helpfully show the
		// request method in OnSend for responses.
		var (
			mu         sync.Mutex
			reqMethods = map[jsonrpc2.ID]string{}
		)
		jsonrpc2.OnRecv(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			switch {
			case req != nil:
				mu.Lock()
				reqMethods[req.ID] = req.Method
				mu.Unlock()

				params, _ := json.Marshal(req.Params)
				if req.Notif {
					logger.Debugf("")
					logger.Printf("jsonrpc2: --> notif: %s: %s\n", req.Method, params)
				} else {
					logger.Printf("jsonrpc2: --> request #%s: %s: %s\n", req.ID, req.Method, params)
				}

			case resp != nil:
				var method string
				if req != nil {
					method = req.Method
				} else {
					method = "(no matching request)"
				}
				switch {
				case resp.Result != nil:
					result, _ := json.Marshal(resp.Result)
					logger.Printf("jsonrpc2: --> result #%s: %s: %s\n", resp.ID, method, result)
				case resp.Error != nil:
					err, _ := json.Marshal(resp.Error)
					logger.Printf("jsonrpc2: --> error #%s: %s: %s\n", resp.ID, method, err)
				}
			}
		})(c)
		jsonrpc2.OnSend(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			switch {
			case req != nil:
				params, _ := json.Marshal(req.Params)
				if req.Notif {
					logger.Printf("jsonrpc2: <-- notif: %s: %s\n", req.Method, params)
				} else {
					logger.Printf("jsonrpc2: <-- request #%s: %s: %s\n", req.ID, req.Method, params)
				}

			case resp != nil:
				mu.Lock()
				method := reqMethods[resp.ID]
				delete(reqMethods, resp.ID)
				mu.Unlock()
				if method == "" {
					method = "(no previous request)"
				}

				if resp.Result != nil {
					result, _ := json.Marshal(resp.Result)
					logger.Printf("jsonrpc2: <-- result #%s: %s: %s\n", resp.ID, method, result)
				} else {
					err, _ := json.Marshal(resp.Error)
					logger.Printf("jsonrpc2: <-- error #%s: %s: %s\n", resp.ID, method, err)
				}
			}
		})(c)
	}
}
//*/
//stdio
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
