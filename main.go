package main

import (
	"fmt"
	"lualsp/auxiliary"
	"lualsp/syntax"
	"os"
	"reflect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("无效参数")
		return
	}
	f, err := os.Open(os.Args[1])
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
		println("line:", line, "col:", col, "msg:", msg)
	})
	lex.Parse()
	//	b, _ := json.MarshalIndent(lex.Block, "", "    ")
	//	fmt.Println(string(b))
	//	fmt.Println("---------------------")
	auxiliary.DrawTree(lex.Block, "All")
	fmt.Println("---------------------")
	auxiliary.Traversal(lex, func(n syntax.Node) {
		fmt.Println(reflect.TypeOf(n).String())
	})
}
