package analysis

import "lualsp/syntax"

func (a *Analysis) analysisAssignStmt(st *syntax.AssignStmt) {
	var rightTypes [][]TypeInfo //第一维索引,第二维类型
	for index, expr := range st.Right {
		switch res := a.analysisExpr(expr).(type) {
		case TypeInfo: //某个类型
			rightTypes = append(rightTypes, []TypeInfo{res})
		case [][]TypeInfo: //某个类型
			if index == len(st.Right)-1 { //如果多类型返回时最后一个则正常添加
				rightTypes = append(rightTypes, res...)
			} else {
				rightTypes = append(rightTypes, res[0]) //否则只添加第一个
			}
		case []*SymbolInfo: //某些对象,getitem返回值
			var types []TypeInfo
			for _, sybif := range res {
				types = append(types, sybif.CurType...)
			}
			rightTypes = append(rightTypes, types)
		case *Symbol: //名字
			if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
				rightTypes = append(rightTypes, sybif.CurType)
			} else {
				rightTypes = append(rightTypes, []TypeInfo{&TypeAny{}})
				//errrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
			}
		case nil:
			rightTypes = append(rightTypes, []TypeInfo{&TypeNil{}})
		default:
			//errrrrrrrrrrrrrrrrrrrrrr
		}
	}

	for index, expr := range st.Left {
		if index >= len(rightTypes) { //右边值不够了
			//errrrrrrrrrrrrrr
			break
		}
		switch res := expr.(type) {
		//如果需要索引
		case *syntax.GetItemExpr:
			if syifs := a.analysisGetItemExpr(res); syifs != nil {
				for _, syif := range syifs {
					syif.CurType = append(syif.CurType, rightTypes[index]...)
				}
			} else {
				//errrrrrrrrrrrrrrrrrrrrrrrrrrrrr
			}
		//是名字
		case *syntax.NameExpr:
			name := a.analysisNameExpr(res)
			if info := a.file.Symbolcur.FindSymbol(name.Name); info != nil {
				info.CurType = append(info.CurType, rightTypes[index]...)
				info.Definitions = append(info.Definitions, name)
			} else {
				name.Types = append(name.Types, rightTypes[index]...)
				//添加到全局变量
				pro := a.file.Project
				syif := &SymbolInfo{
					Definitions: []*Symbol{name},
					References:  []*Symbol{name},
				}
				syif.CurType = append(syif.CurType, name.Types...)
				pro.SymbolsMu.Lock()
				pro.SymbolList.Symbols[name.Name] = syif
				pro.SymbolsMu.Unlock()
			}
		default:
			//errrrrrrrrrrrrrrrrrrrrrr
		}
	}

}
func (a *Analysis) analysisLabelStmt(st *syntax.LabelStmt) {
	if nameExpr, ok := st.Name.(*syntax.NameExpr); ok {
		name := a.analysisNameExpr(nameExpr)
		if syif := a.file.Symbolcur.FindLabel(name.Name); syif != nil {
			if len(syif.Definitions) != 0 { //正常
				//errrrrrrrrrrrrr重复定义
				return
			}
		}

		syif := &SymbolInfo{
			CurType:     []TypeInfo{&Typelabel{Value: name.Name}},
			Definitions: []*Symbol{name},
		}
		a.file.Symbolcur.Labels[name.Name] = syif
		lists := a.file.Symbolcur.FindlonelyLabel(name.Name) //查找内部label
		for _, list := range lists {
			info := list.Labels[name.Name]
			if len(info.Definitions) == 0 { //如果是个空虚的label
				syif.References = append(syif.References, info.References...) //上车
				delete(list.Labels, name.Name)                                //释放
			}
		}
		a.file.Symbolcur.Labels[name.Name] = syif
	} else {
		//errrrrrrrrrrrrrr
	}
}
func (a *Analysis) analysisBreakStmt(st *syntax.BreakStmt) {
	var temp *SymbolList
	//循环向外层寻找
	for temp = a.file.Symbolcur; temp != nil; temp = temp.Outside {
		switch temp.Node.(type) {
		case *syntax.ForLoopListStmt, *syntax.ForLoopNumStmt, *syntax.WhileStmt:
			break
		}
	}
	if temp == nil {
		//errrrrrrrrrrrrrrrrrrr
	}
}
func (a *Analysis) analysisGotoStmt(st *syntax.GotoStmt) {
	if nameExpr, ok := st.Name.(*syntax.NameExpr); ok {
		name := a.analysisNameExpr(nameExpr)
		if syif := a.file.Symbolcur.FindLabel(name.Name); syif != nil {
			syif.References = append(syif.References, name)
			name.SymbolInfo = syif
		} else {
			syif := &SymbolInfo{
				CurType:    []TypeInfo{&Typelabel{Value: name.Name}},
				References: []*Symbol{name},
			}
			name.SymbolInfo = syif
			a.file.Symbolcur.Labels[name.Name] = syif //空头作用域
		}
	} else {
		//errrrrrrrrrrrrrr
	}
}
func (a *Analysis) analysisDoEndStmt(st *syntax.DoEndStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
}
func (a *Analysis) analysisWhileStmt(st *syntax.WhileStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	switch res := a.analysisExpr(st.Condition).(type) {
	case *Symbol: //名字
		if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
			res.SymbolInfo = sybif
			sybif.References = append(sybif.References, res)
		} else {
			//errrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
		}
	case nil:
		//errrrrrrrrrrrrrrrrrrrrrr
	default:
		//errrrrrrrrrrrrrrrrrrrrrr
	}
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
}
func (a *Analysis) analysisRepeatStmt(st *syntax.RepeatStmt) {
}
func (a *Analysis) analysisIfStmt(st *syntax.IfStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	switch res := a.analysisExpr(st.Condition).(type) {
	case *Symbol: //名字
		if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
			res.SymbolInfo = sybif
			sybif.References = append(sybif.References, res)
		} else {
			//errrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
		}
	case nil:
		//errrrrrrrrrrrrrrrrrrrrrr
	default:
		//errrrrrrrrrrrrrrrrrrrrrr
	}
	for _, stmt := range st.Then {
		a.analysisStmt(stmt)
	}
	for _, stmt := range st.Else {
		a.analysisStmt(stmt)
	}
}
func (a *Analysis) analysisForLoopNumStmt(st *syntax.ForLoopNumStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	//分析名称
	if nameExpr, ok := st.Name.(*syntax.NameExpr); ok {
		res := a.analysisNameExpr(nameExpr)
		res.SymbolInfo = &SymbolInfo{
			CurType:     []TypeInfo{&TypeNumber{}},
			Definitions: []*Symbol{res},
			References:  []*Symbol{res},
		}
		a.file.Symbolcur.Symbols[res.Name] = res.SymbolInfo
	} else {

	}
	//分析后面的数字表达式
	for _, exp := range []syntax.Expr{st.Init, st.Limit, st.Step} {
		switch res := a.analysisExpr(exp).(type) {
		case *Symbol: //名字
			if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
				res.SymbolInfo = sybif
				sybif.References = append(sybif.References, res)
				//检查类型是否为数字
			} else {
				//errrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
			}
		case [][]TypeInfo: //检查类型是否为数字
		case TypeInfo: //检查类型是否为数字
		case []*SymbolInfo: //检查类型是否为数字
		case nil:
			//errrrrrrrrrrrrrrrrrrrrrr
		default:
			//errrrrrrrrrrrrrrrrrrrrrr
		}
	}

	//分析内容
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
}
func (a *Analysis) analysisForLoopListStmt(st *syntax.ForLoopListStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	//分析名称
	for _, name := range st.Names {
		if nameExpr, ok := name.(*syntax.NameExpr); ok {
			res := a.analysisNameExpr(nameExpr)
			res.SymbolInfo = &SymbolInfo{
				CurType:     []TypeInfo{&TypeNumber{}},
				Definitions: []*Symbol{res},
				References:  []*Symbol{res},
			}
			a.file.Symbolcur.Symbols[res.Name] = res.SymbolInfo
		} else {

		}
	}
	//分析迭代表达式
	switch expr := st.Exprs[0].(type) {
	case *syntax.FuncCall:
		res := a.analysisFuncCall(expr) //返回值的第一个必须是function
		if len(res) > 0 {
			count := 0
			for _, tp := range res[0] {
				if tp.TypeName() == "function" {
					count++
				}
			}
			if count == 0 {
				//errrrrrrrrrrrrrr
			}
		} else {
			//errrrrrrrrrrrrrr
		}
	case *syntax.NameExpr:

	}

}
func (a *Analysis) analysisFuncDefStmt(st *syntax.FuncDefStmt) {
}
func (a *Analysis) analysisLocalFuncDefStmt(st *syntax.LocalFuncDefStmt) {
}
func (a *Analysis) analysisLocalVarDef(st *syntax.LocalVarDef) {
}
func (a *Analysis) analysisReturnStmt(st *syntax.ReturnStmt) {
}
func (a *Analysis) analysisErrorStmt(st *syntax.ErrorStmt) {
}
