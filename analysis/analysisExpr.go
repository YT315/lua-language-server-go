package analysis

import (
	"lualsp/syntax"
	"regexp"
	"strconv"
	"strings"
)

type addType func(*Symbol)

//正则表达式用于判断类型信息的格式
var (
	setTypeReg   = regexp.MustCompile("^--\\[\\[@(?P<types>[A-Za-z0-9_]+(\\|[A-Za-z0-9]+)*)\\]\\]$")
	setTypeIndex = setTypeReg.SubexpIndex("types")
	addTypeReg   = regexp.MustCompile("^--\\[\\[+(?P<type>[A-Za-z0-9_]+)\\]\\]$")
	addTypeIndex = addTypeReg.SubexpIndex("type")
)

func (a *Analysis) analysisNameExpr(ep *syntax.NameExpr) (result *Symbol) {
	result = &Symbol{
		Name: ep.Value,
		Node: ep,
		File: a.file,
	}
	switch data := ep.Type.(type) {
	case *syntax.ATypeExpr:
		if res := a.analysisATypeExpr(data); res != nil {
			res(result)
		}
	case *syntax.STypeExpr:
		if res := a.analysisSTypeExpr(data); res != nil {
			result.Types = append(result.Types, res...)
		}
	case nil:
	default:
		//errrrrrrrrrrrrrrrr
	}
	return result
}

//type
//设置类型@
func (a *Analysis) analysisSTypeExpr(ep *syntax.STypeExpr) (result []TypeInfo) {
	result = []TypeInfo{}
	if temp := setTypeReg.FindStringSubmatch(ep.Value); len(temp) > 0 {
		types := strings.Split(temp[setTypeIndex], "|")
		a.file.Project.TypesMu.RLock()
		for _, typename := range types {
			if typeSyb, okay := a.file.Project.TypeList[typename]; okay {
				result = append(result, typeSyb.Types...)
			} else {
				//errrrrrrrrrrrrrrrrrrrr
			}
		}
		a.file.Project.TypesMu.RUnlock()
	} else {
		//errrrrrrrrrrrrrrrrrrrrrrrr
	}
	return result
}

//添加类型+
func (a *Analysis) analysisATypeExpr(ep *syntax.ATypeExpr) (result addType) {
	if temp := addTypeReg.FindStringSubmatch(ep.Value); len(temp) > 0 {
		typename := temp[addTypeIndex]
		a.file.Project.TypesMu.RLock()
		if _, okay := a.file.Project.TypeList[typename]; okay {
			//warrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
		}
		a.file.Project.TypesMu.RUnlock()
		result = func(syb *Symbol) {
			a.file.Project.TypesMu.Lock()
			a.file.Project.TypeList[typename] = syb
			a.file.Project.TypesMu.Unlock()
		}
	} else {
		//errrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
	}
	return result
}

//TypeInfo
func (a *Analysis) analysisNilExpr(ep *syntax.NilExpr) TypeInfo {
	return nil
}
func (a *Analysis) analysisFalseExpr(ep *syntax.FalseExpr) TypeInfo {
	return &TypeBool{Value: false}
}
func (a *Analysis) analysisTrueExpr(ep *syntax.TrueExpr) TypeInfo {
	return &TypeBool{Value: true}
}
func (a *Analysis) analysisNumberExpr(ep *syntax.NumberExpr) TypeInfo {
	num, _ := strconv.ParseFloat(ep.Value, 64)
	return &TypeNumber{Value: num}
}
func (a *Analysis) analysisStringExpr(ep *syntax.StringExpr) TypeInfo {
	return &TypeString{Value: ep.Value}
}
func (a *Analysis) analysisAnyExpr(ep *syntax.AnyExpr) TypeInfo {
	return &TypeAny{}
}

//解析函数体
func (a *Analysis) analysisFuncDefExpr(ep *syntax.FuncDefExpr) (result *TypeFunction) {
	a.file.createInside()      //创建新作用域
	defer a.file.backOutside() //退出作用域
	//分析参数
	if pe, ok := ep.Param.(*syntax.ParamExpr); ok {
		sybs := a.analysisParamExpr(pe)
		for _, syb := range sybs {
			syb.SymbolInfo = &SymbolInfo{
				Definitions: []*Symbol{syb},
				References:  []*Symbol{syb},
			}
			syb.SymbolInfo.CurType = append(syb.SymbolInfo.CurType, syb.Types...) //添加类型
			a.file.Symbolcur.Symbols[syb.Name] = syb.SymbolInfo
		}
	} else {
		//errrrrrrrrrrrrrrrrrrrrrrr
	}
	//分析语句
	for _, stmt := range ep.Block {
		a.analysisStmt(stmt)
	}
	//分析返回值
	if ste, ok := ep.Param.(*syntax.STypeExpr); ok {
		result = &TypeFunction{
			Returns: a.analysisSTypeExpr(ste),
		}
	} else {
		//errrrrrrrrrrrr
	}
	return
}

//解析函数参数
func (a *Analysis) analysisParamExpr(ep *syntax.ParamExpr) (result []*Symbol) {
	result = make([]*Symbol, len(ep.Params))
	for _, param := range ep.Params {
		if nx, ok := param.(*syntax.NameExpr); ok {
			if syb := a.analysisNameExpr(nx); syb != nil { //解析成功
				result = append(result, syb)
			} else {
				//errrrrrrrrrrrrrrrrrrrrr
			}
		} else {
			//errrrrrrrrrrrrrrrrrrrrrrrr
		}
	}
	return result
}

//解析属性获取
func (a *Analysis) analysisGetItemExpr(ep *syntax.GetItemExpr) (result *Symbol) {
	var tab *Symbol
	switch data := ep.Table.(type) {
	case *syntax.GetItemExpr:
		tab = a.analysisGetItemExpr(data)
	case *syntax.NameExpr:
		//寻找此符号,查看符号类型,如果不是table,//报警并添加table类型,是table则ok
		name := a.analysisNameExpr(data)
		if info := a.file.Symbolcur.FindSymbol(name.Name); info != nil {

		}
	case *syntax.FuncCall:
	default:
		//errrrrrrrrrrrrrrrr
	}
	return nil
}
func (a *Analysis) analysisTableExpr(ep *syntax.TableExpr) interface{} {
	return nil
}
func (a *Analysis) analysisFieldExpr(ep *syntax.FieldExpr) interface{} {
	return nil
}
func (a *Analysis) analysisTwoOpExpr(ep *syntax.TwoOpExpr) interface{} {
	return nil
}
func (a *Analysis) analysisOneOpExpr(ep *syntax.OneOpExpr) interface{} {
	return nil
}
