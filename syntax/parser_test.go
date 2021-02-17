package syntax

import (
	"fmt"
	"lualsp/auxiliary"
	"os"
	"testing"
)

func Test_GetWriter(t *testing.T) {

	f, err := os.Open("../test.lua")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	lex := NewLexer(f, func(line, col uint, msg string) {
		println("line:", line, "col:", col, "msg:", msg)
	})
	lex.Parse()
	auxiliary.DrawTree(lex.Block, "All")
	t.Log("ok")
}
