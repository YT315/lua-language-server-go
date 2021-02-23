package analysis

import "lualsp/syntax"

//Project 工程对象
type Project struct {
	Workspaces []*Workspace //工作区列表
}

//Workspace 工作区
type Workspace struct {
	RootPath string
	Files    map[string]*File
}

//File 文件对象
type File struct {
	Name       string                //文件名,不包含后缀
	Path       string                //文件路径,包括文件名
	content    string                //文件内容实时
	IsEditing  bool                  //文件是否正在编辑,编辑时,文件的实际内容和content不一定相同
	Ast        []syntax.Stmt         //文件的抽象语法树
	SymbolPos  []*Symbol             //文件中所有符号列表,按照位置顺序向后排列
	SymbolList *SymbolList           //文件符号表,作用域
	TypeList   map[string]SymbolInfo //文件中包含的所有类型列表
	BeRequire  []*File               //所有依赖此文件的文件
}

//SymbolList 符号表结构
type SymbolList struct {
	Position int                //此符号表作用范围
	Outside  *SymbolList        //此符号表外部符号表
	Inside   []*SymbolList      //此符号表内部符号表
	Symbols  map[string]*Symbol //包含的符号
}

//Symbol 符号对象
type Symbol struct {
	Node       syntax.Node //此符号对应的语法树节点
	File       *File       //此符号所在文件
	TypeInfo   TypeInfo    //符号的类型
	SymbolInfo *SymbolInfo //符号的符号信息,引用以及定义
}

//SymbolInfo 符号信息
type SymbolInfo struct {
	Definitions []*Symbol //符号定义处
	References  []*Symbol //符号所有引用处
}

//TypeInfo 类型接口
type TypeInfo interface {
	TypeName() string //类型名称
}

//TypeBool 布尔类型
type TypeBool struct{}

//TypeName 类型名称
func (*TypeBool) TypeName() string {
	return "bool"
}

//TypeNumber 数字类型
type TypeNumber struct{}

//TypeName 类型名称
func (*TypeNumber) TypeName() string {
	return "number"
}

//TypeString 字符串类型
type TypeString struct{}

//TypeName 类型名称
func (*TypeString) TypeName() string {
	return "string"
}

//TypeAny 任意类型
type TypeAny struct{}

//TypeName 类型名称
func (*TypeAny) TypeName() string {
	return "any"
}

//TypeTable 表类型
type TypeTable struct {
	Name        string
	IsAnonymous bool
	Fields      map[string]TypeInfo
}

//TypeName 类型名称
func (me *TypeTable) TypeName() string {
	return me.Name
}

//TypeFunction 函数类型
type TypeFunction struct {
	Name     string
	Returns  []TypeInfo
	Symbols  map[string]*Symbol
	TypeList []TypeInfo
}

//TypeName 类型名称
func (me *TypeFunction) TypeName() string {
	return me.Name
}
