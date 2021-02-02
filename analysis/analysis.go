package analysis

import (
	"lualsp/syntax"
)

type Project struct {
	files map[string]syntax.Node //文件分析的语法树
	scope string                 //全局变量

}
