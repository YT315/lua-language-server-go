package analysis

import (
	"bytes"
	"io/ioutil"
	"lualsp/protocol"
	"lualsp/syntax"
	"sync"
)

//File 文件对象
type File struct {
	Project     *Project      //所属工程
	Name        string        //文件名,不包含后缀
	Path        string        //文件路径,包括文件名
	content     []byte        //文件内容实时
	linePos     map[int]int   //行号对应的字节偏移方便有变动时插入
	IsEditing   bool          //文件是否正在编辑,编辑时,文件的实际内容和content不一定相同
	Ast         []syntax.Stmt //文件的抽象语法树
	Diagnostics []protocol.Diagnostic
	SymbolPos   []*Symbol    //文件中所有符号列表,按照位置顺序向后排列
	Symbolbase  *SymbolList  //文件符号表,作用域深度为1层
	Symbolcur   *SymbolList  //文件符号表,作用域
	ReturnType  [][]TypeInfo //返回列表
	BeRequire   []*File      //所有依赖此文件的文件
	Mutex       sync.Mutex   //互斥锁
}

//创建一个新的作用域
func (f *File) createInside(node syntax.Node) (list *SymbolList) {
	list = &SymbolList{
		Node:    node,
		Deep:    f.Symbolcur.Deep + 1,
		Outside: f.Symbolcur,
	}
	f.Symbolcur.Inside = append(f.Symbolcur.Inside, list)
	f.Symbolcur = list
	return
}

//返回上一层作用域
func (f *File) backOutside() (list *SymbolList) {
	if f.Symbolcur.Outside != nil {
		f.Symbolcur = f.Symbolcur.Outside
		list = f.Symbolcur
	}
	return
}

//将文件内容读取缓存,并更新行号表
func (f *File) updata() (err error) {
	f.content, err = ioutil.ReadFile(f.Path)

	return
}

//解析
func (f *File) Parse() {
	reader := bytes.NewReader(f.content)
	lex := syntax.NewLexer(reader, func(line, col uint, msg string) {
		println("err:- line:", line, "col:", col, "msg:", msg)
	})
	lex.Parse()
	f.Ast = append(f.Ast, lex.Block...)
	f.Diagnostics = append(f.Diagnostics, lex.Diagnostics...)
}

//文件内容
type FileContent struct {
	content []byte      //文件内容实时
	linePos map[int]int //行号对应的字节偏移方便有变动时插入
}

//	重写内容
func (f *FileContent) Overwrite(content []byte) {
	f.content = nil //释放原来内存
	f.linePos = map[int]int{}
	f.content = append(f.content, content...) //加入内容
	line := 1
	f.linePos[line] = 0
	for index, value := range f.content { //更新行号
		if value == '\n' {
			line++
			f.linePos[line] = index
		}
	}
}

//	重写内容
func (f *FileContent) Insert(starline, staroff, endline, endoff, RangeLengthint int, text string) {
	newtext := []byte(text)
	starIndex := f.linePos[starline] + staroff
	endIndex := f.linePos[endline] + endoff

	f.content = append(append(f.content[:starIndex], newtext...), f.content[endIndex:]...)

}
