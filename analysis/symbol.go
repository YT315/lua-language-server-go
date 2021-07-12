package analysis

import "lualsp/syntax"

//SymbolList 符号表结构
type SymbolList struct {
	Deep    int                    //此作用域的深度,全局变量深度为0,文件为1
	Node    syntax.Node            //此范围对应的语法树节点,if/while,for....
	Outside *SymbolList            //此符号表外部符号表
	Inside  []*SymbolList          //此符号表内部符号表
	Symbols map[string]*SymbolInfo //包含的符号
	Labels  map[string]*SymbolInfo //包含的标签
}

//查找符号
func (sl *SymbolList) FindSymbol(name string) (result *SymbolInfo) {
	var temp *SymbolList
	//循环向外层寻找
	for temp = sl; temp != nil; temp = temp.Outside {
		if info, ok := temp.Symbols[name]; ok {
			result = info
			break
		}
	}
	return
}

//查找label
func (sl *SymbolList) FindLabel(name string) (result *SymbolInfo) {
	var temp *SymbolList
	//循环向外层寻找
	for temp = sl; temp != nil; temp = temp.Outside {
		if info, ok := temp.Labels[name]; ok {
			result = info
			break
		}
		if _, ok := temp.Node.(*syntax.FuncDefExpr); ok { //label仅函数级寻找
			break
		}
	}
	return
}

//查找空的label
func (sl *SymbolList) FindlonelyLabel(name string) (list []*SymbolList) {
	if _, ok := sl.Symbols[name]; ok {
		list = append(list, sl)
	}
	//循环向内层寻找
	for _, temp := range sl.Inside {
		if _, ok := temp.Node.(*syntax.FuncDefExpr); !ok { //函数级寻找
			list = append(list, temp.FindlonelyLabel(name)...)
		}
	}
	return
}

//Symbol 符号对象
//代表代码中某个位置的唯一一个符号,和源码中的文本一一对应
type Symbol struct {
	Name      string
	Node      syntax.Node //此符号对应的语法树节点
	File      *File       //此符号所在文件
	Types     *TypeSet    //符号的类型在未知判断的情况下可能有多个
	SymbolCtx *SymbolInfo //符号的上下文信息,引用以及定义
}

//SymbolInfo 符号信息
//当一个符号有上下文时,此对象保存此符号的上下文信息
type SymbolInfo struct {
	CurType     *TypeSet  //符号当前的类型
	Definitions []*Symbol //符号定义处
	References  []*Symbol //符号所有引用处
}
