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
	StateMu    sync.Mutex
	State      projectState
	Wg         sync.WaitGroup //等待工作区初始化完成
	Workspaces []*Workspace   //工作区列表
	SymbolList *SymbolList    //全局符号表
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
					Name:       info.Name(),
					Path:       path,
					Workspace:  w,
					SymbolList: w.Project.SymbolList,
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
	Workspace  *Workspace
	Name       string                //文件名,不包含后缀
	Path       string                //文件路径,包括文件名
	content    []byte                //文件内容实时
	IsEditing  bool                  //文件是否正在编辑,编辑时,文件的实际内容和content不一定相同
	Ast        []syntax.Stmt         //文件的抽象语法树
	SymbolPos  []*Symbol             //文件中所有符号列表,按照位置顺序向后排列
	SymbolList *SymbolList           //文件符号表,作用域
	TypeList   map[string]SymbolInfo //文件中包含的所有类型列表
	BeRequire  []*File               //所有依赖此文件的文件
	Mutex      sync.Mutex            //互斥锁
}

func (f *File) updata() (err error) {
	f.content, err = ioutil.ReadFile(f.Path)
	return
}
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
	Deep     int                //此作用域的深度
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
