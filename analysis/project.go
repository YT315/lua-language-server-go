package analysis

import (
	"context"
	"lualsp/logger"
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
func (p *Project) Scan(ctx context.Context) {
	p.StateMu.Lock()
	p.State = ProjectFileScanning
	p.StateMu.Unlock()
	for _, ws := range p.Workspaces {
		p.Wg.Add(1)
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
				path = filepath.ToSlash(path)
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
