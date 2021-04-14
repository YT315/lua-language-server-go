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
	//分析语句////////////////////////////////////////////获取中间返回值
	for _, stmt := range ep.Block {
		a.analysisStmt(stmt)
	}
	//分析返回值
	if ep.Result != nil {
		if ste, ok := ep.Result.(*syntax.STypeExpr); ok {
			result = &TypeFunction{}
			result.Returns = append(result.Returns, a.analysisSTypeExpr(ste))
		} else {
			//errrrrrrrrrrrr
		}
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
func (a *Analysis) analysisGetItemExpr(ep *syntax.GetItemExpr) (result []*SymbolInfo) {
	//------------获取tab表达式的类型
	var types []*TypeTable
	switch data := ep.Table.(type) {
	//如果需要索引
	case *syntax.GetItemExpr:
		if syifs := a.analysisGetItemExpr(data); syifs != nil {
			for _, syif := range syifs {
				for _, tp := range syif.CurType {
					if tp.TypeName() == "table" {
						types = append(types, tp) //将table类型添加
					}
				}
			}
		} else {
			//errrrrrrrrrrrrrrrrrrrrrrrrrrrrr
			return nil
		}
	//是名字
	case *syntax.NameExpr:
		//寻找此符号,查看符号类型,如果不是table,//报警并添加table类型,是table则ok
		name := a.analysisNameExpr(data)
		if info := a.file.Symbolcur.FindSymbol(name.Name); info != nil {
			info.CurType = append(info.CurType, name.Types...)
			//判断是否有表类型
			for _, tp := range info.CurType {
				if tp.TypeName() == "table" {
					types = append(types, tp) //将table类型添加
				}
			}
			//未找到表类型,则给name添加表类型,并报错误
			if len(types) == 0 {
				//errrrrrrrrrrrrrrrrrrr
				//创建一个表类型
				newtab := &TypeTable{}
				name.SymbolInfo.CurType = append(name.SymbolInfo.CurType, newtab)
				types = append(types, newtab) //添加到返回值
			}
		} else {
			//errrrrrrrrrrrrrrrrrrr
		}
	case *syntax.FuncCall:
		funres := a.analysisFuncCall(data)
		if len(funres.Returns) > 0 {
			for _, tp := range funres.Returns[0] {
				if tp.TypeName() == "table" {
					types = append(types, tp) //将table类型添加
				}
			}
		} else {
			//errrrrrrrrrrrrrrrrrrrrrr
			return nil
		}
	default:
		//errrrrrrrrrrrrrrrr
		return nil
	}
	if len(types) == 0 {
		return nil
	}
	//------------获取key表达式的类型
	switch data := ep.Key.(type) {
	//是字符串
	case *syntax.StringExpr:
		for _,ty:=range types{
			if 
		}
	//是数字
	case *syntax.NumberExpr:
	default:
		//errrrrrrrrrrrrrrrr
		return nil
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
