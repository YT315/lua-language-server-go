package analysis

import (
	"bytes"
	"io/ioutil"
	"lualsp/logger"
	"lualsp/syntax"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type projectState int

const (
	ProjectCreated      = projectState(iota)
	ProjectFileScanning //正在扫描并计算AST
	ProjectFileScanned  //扫描完成
	ProjectAnalysising  //正在分析
	ProjectAnalysised   //分析完成
)

//Project 工程对象
type Project struct {
	//工程状态机
	StateMu    sync.RWMutex
	State      projectState
	Wg         sync.WaitGroup     //等待工作区初始化完成
	Workspaces []*Workspace       //工作区列表
	SymbolsMu  sync.RWMutex       //符号表读写锁
	SymbolList SymbolList         //全局符号表
	TypesMu    sync.RWMutex       //类型表读写锁
	TypeList   map[string]*Symbol //所有类型列表,类型其实是对某个符号的引用
}

//扫描所有工作区
func (p *Project) Scan() {
	p.StateMu.Lock()
	p.State = ProjectFileScanning
	p.StateMu.Unlock()
	for _, ws := range p.Workspaces {
		go ws.Scan(&p.Wg)
	}
	p.Wg.Wait()
	p.StateMu.Lock()
	p.State = ProjectFileScanned
	p.StateMu.Unlock()
	logger.Debugln("scan finish")
	p.analysis()
}

//开始分析
func (p *Project) analysis() {
	p.StateMu.Lock()
	p.State = ProjectAnalysising
	p.StateMu.Unlock()

	p.StateMu.Lock()
	p.State = ProjectAnalysised
	p.StateMu.Unlock()

}

//Workspace 工作区
type Workspace struct {
	Project  *Project
	RootPath string
	Files    map[string]*File //key:文件的相对路径
}

//扫描所有文件
func (w *Workspace) Scan(wg *sync.WaitGroup) {
	//限制同时打开的文件数量
	sem := make(chan struct{}, runtime.GOMAXPROCS(0)+10)
	wg.Add(1)
	err := filepath.Walk(w.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".lua") { //只要lua文件
				f := &File{
					Name:    info.Name(),
					Path:    path,
					Project: w.Project,
				}
				w.Files[path] = f
				wg.Add(1)
				go func() {
					sem <- struct{}{} //控制并发数量
					defer func() { <-sem }()
					if err := f.updata(); err != nil {
						logger.Warningln(err.Error())
					} else {
						f.Parse()
					}
					wg.Done()
				}()
			}
			return nil
		})
	if err != nil {
		logger.Errorln(err)
	}
	wg.Done()
}

//File 文件对象
type File struct {
	Project    *Project      //所属工程
	Name       string        //文件名,不包含后缀
	Path       string        //文件路径,包括文件名
	content    []byte        //文件内容实时
	linePos    map[int]int   //行号对应的字节偏移方便有变动时插入
	IsEditing  bool          //文件是否正在编辑,编辑时,文件的实际内容和content不一定相同
	Ast        []syntax.Stmt //文件的抽象语法树
	SymbolPos  []*Symbol     //文件中所有符号列表,按照位置顺序向后排列
	Symbolbase *SymbolList   //文件符号表,作用域深度为1层
	Symbolcur  *SymbolList   //文件符号表,作用域
	ReturnType [][]TypeInfo  //返回列表
	BeRequire  []*File       //所有依赖此文件的文件
	Mutex      sync.Mutex    //互斥锁
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

//将文件内容读取缓存
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
	f.Ast = lex.Block
}

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
		if info, ok := temp.Symbols[name]; ok {
			result = info
			break
		}
		if _, ok := temp.Node.(*syntax.FuncDefExpr); ok { //函数级寻找
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
type Symbol struct {
	Name       string
	Node       syntax.Node //此符号对应的语法树节点
	File       *File       //此符号所在文件
	Types      []TypeInfo  //符号的类型在未知判断的情况下可能有多个
	SymbolInfo *SymbolInfo //符号的符号信息,引用以及定义
}

//SymbolInfo 符号信息
type SymbolInfo struct {
	CurType     []TypeInfo //符号当前的类型
	Definitions []*Symbol  //符号定义处
	References  []*Symbol  //符号所有引用处
}

//TypeInfo 类型接口
type TypeInfo interface {
	TypeName() string //类型名称
}

//TypeNil 空类型
type TypeNil struct {
}

//TypeName 类型名称
func (*TypeNil) TypeName() string {
	return "nil"
}

//TypeBool 布尔类型
type TypeBool struct {
	Value bool
}

//TypeName 类型名称
func (*TypeBool) TypeName() string {
	return "bool"
}

//TypeNumber 数字类型
type TypeNumber struct {
	Value float64
}

//TypeName 类型名称
func (*TypeNumber) TypeName() string {
	return "number"
}

//TypeString 字符串类型
type TypeString struct {
	Value string
}

//TypeName 类型名称
func (*TypeString) TypeName() string {
	return "string"
}

//Typelabel 字符串类型
type Typelabel struct {
	Value string
}

//TypeName 类型名称
func (*Typelabel) TypeName() string {
	return "label"
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
	Fields      map[string]*SymbolInfo  //hash
	Items       map[float64]*SymbolInfo //array
	Metatable   *TypeTable              //元表
}

//TypeName 类型名称
func (me *TypeTable) TypeName() string {
	return "table"
}

//TypeFunction 函数类型
type TypeFunction struct {
	Returns [][]TypeInfo //函数可能有多种返回值情况,数字第一索引表示,返回值索引,第二索引,表示此返回值的类型范围
}

//TypeName 类型名称
func (me *TypeFunction) TypeName() string {
	return "function"
}
