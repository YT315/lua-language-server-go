package analysis

import "lualsp/syntax"

//赋值
func (a *Analysis) analysisAssignStmt(st *syntax.AssignStmt) {
	var rightTypes [][]TypeInfo //第一维索引,第二维类型
	//分析右边
	for index, expr := range st.Right {
		switch res := a.analysisExpr(expr).(type) {
		case TypeInfo: //某个直接的类型
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
				types = append(types, sybif.CurType.Types...)
			}
			rightTypes = append(rightTypes, types)
		case *Symbol: //名字
			if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
				rightTypes = append(rightTypes, sybif.CurType.Types)
			} else {
				rightTypes = append(rightTypes, []TypeInfo{typeAny})
				//未找到定义处,先放入全局变量中最后在遍历全局变量,确定未定义的变量
				pro := a.file.Project
				syif := &SymbolInfo{
					References: []*Symbol{res},
					CurType:    NewTypeSetWithContent(typeAny),
				}
				pro.SymbolsMu.Lock()
				pro.SymbolList.Symbols[res.Name] = syif
				pro.SymbolsMu.Unlock()
			}
		case nil:
			rightTypes = append(rightTypes, []TypeInfo{typeNil})
		default:
			//类型错误,不能作为右值
			rightTypes = append(rightTypes, []TypeInfo{typeAny})
			err := &AnalysisErr{Errtype: NotRightValue}
			err.Scope = expr.GetScope()
			err.insertInto(a)
		}
	}
	//不存在无左值错误时
	if st.Err == nil || st.Err.Errtype != syntax.LackLeft {
		for index, expr := range st.Left {
			types := []TypeInfo{typeAny}
			if index < len(rightTypes) {
				types = rightTypes[index]
			}
			switch res := expr.(type) {
			//是名字
			case *syntax.NameExpr:
				name := a.analysisNameExpr(res)
				if info := a.file.Symbolcur.FindSymbol(name.Name); info != nil {
					info.CurType.AddRange(types...)
					info.Definitions = append(info.Definitions, name)
				} else {
					name.Types.AddRange(types...)
					//添加到全局变量
					pro := a.file.Project
					syif := &SymbolInfo{
						Definitions: []*Symbol{name},
						References:  []*Symbol{name},
						CurType:     NewTypeSetWithContent(name.Types.Types...),
					}
					pro.SymbolsMu.Lock()
					pro.SymbolList.Symbols[name.Name] = syif
					pro.SymbolsMu.Unlock()
				}
			//如果需要索引
			case *syntax.GetItemExpr:
				if syifs := a.analysisGetItemExpr(res); syifs != nil {
					for _, syif := range syifs {
						syif.CurType.AddRange(types...)
					}
				} else {
					err := &AnalysisErr{Errtype: IndexErr}
					err.Scope = expr.GetScope()
					err.insertInto(a)
				}
			default:
				err := &AnalysisErr{Errtype: TypeErr}
				err.Scope = expr.GetScope()
				err.insertInto(a)
			}
		}
	}

}

//标签
func (a *Analysis) analysisLabelStmt(st *syntax.LabelStmt) {
	if nameExpr, ok := st.Name.(*syntax.NameExpr); ok {
		//标签不需要格式
		if nameExpr.Type != nil {
			err := &AnalysisErr{Errtype: LabelFormatErr}
			err.Scope = nameExpr.GetScope()
			err.insertInto(a)
		}
		name := a.analysisNameExpr(nameExpr)
		//类型设置为lab
		labty := &Typelabel{Value: name.Name}
		name.Types.Add(labty)

		//先向外层寻找重复的lab
		if syif := a.file.Symbolcur.FindLabel(name.Name); syif != nil {
			if len(syif.Definitions) != 0 { //重复定义
				err := &AnalysisErr{Errtype: LabelRedef}
				err.Scope = st.GetScope()
				err.insertInto(a)
			}
			//如果找到的目标和当前是同级,则加入
			if info, ok := a.file.Symbolcur.Labels[name.Name]; ok && info == syif {
				syif.Definitions = append(syif.Definitions, name)
				name.SymbolCtx = syif
			}
		} else {
			syif := &SymbolInfo{
				CurType:     NewTypeSetWithContent(name.Types.Types...),
				Definitions: []*Symbol{name},
			}
			lists := a.file.Symbolcur.FindlonelyLabel(name.Name) //查找内部label
			for _, list := range lists {
				info := list.Labels[name.Name]
				if len(info.Definitions) == 0 { //如果是个空虚的label
					syif.References = append(syif.References, info.References...) //上车
					for _, syb := range info.References {
						syb.SymbolCtx = syif //上下文
					}
					delete(list.Labels, name.Name) //释放
				}
			}
			a.file.Symbolcur.Labels[name.Name] = syif
		}
	} else if st.Name != nil { //语法解析错误
		err := &AnalysisErr{Errtype: SyntaxDataErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
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
		err := &AnalysisErr{Errtype: BreakNoLoopErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
	}
}

func (a *Analysis) analysisGotoStmt(st *syntax.GotoStmt) {
	if nameExpr, ok := st.Name.(*syntax.NameExpr); ok {
		name := a.analysisNameExpr(nameExpr)
		//类型设置为lab
		labty := &Typelabel{Value: name.Name}
		name.Types.Add(labty)

		if syif := a.file.Symbolcur.FindLabel(name.Name); syif != nil {
			syif.References = append(syif.References, name)
			name.SymbolCtx = syif
		} else {
			syif := &SymbolInfo{
				CurType:    NewTypeSetWithContent(name.Types.Types...),
				References: []*Symbol{name},
			}
			name.SymbolCtx = syif
			a.file.Symbolcur.Labels[name.Name] = syif //空头作用域
		}
	} else if st.Name != nil { //语法解析错误
		err := &AnalysisErr{Errtype: SyntaxDataErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
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
	if res, okay := a.analysisExpr(st.Condition).(*Symbol); okay {
		if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
			res.SymbolCtx = sybif
			sybif.References = append(sybif.References, res)
		} else {
			err := &AnalysisErr{Errtype: NoDefine}
			err.Scope = st.GetScope()
			err.insertInto(a)
		}
	}
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
}

func (a *Analysis) analysisRepeatStmt(st *syntax.RepeatStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
	if res, okay := a.analysisExpr(st.Condition).(*Symbol); okay {
		if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
			res.SymbolCtx = sybif
			sybif.References = append(sybif.References, res)
		} else {
			err := &AnalysisErr{Errtype: NoDefine}
			err.Scope = st.GetScope()
			err.insertInto(a)
		}
	}

}

func (a *Analysis) analysisIfStmt(st *syntax.IfStmt) {
	a.file.createInside(st)    //创建新作用域
	defer a.file.backOutside() //退出作用域
	if res, okay := a.analysisExpr(st.Condition).(*Symbol); okay {
		if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
			res.SymbolCtx = sybif
			sybif.References = append(sybif.References, res)
		} else {
			err := &AnalysisErr{Errtype: NoDefine}
			err.Scope = st.GetScope()
			err.insertInto(a)
		}
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
		res.SymbolCtx = &SymbolInfo{
			CurType:     NewTypeSetWithContent(typeNumber),
			Definitions: []*Symbol{res},
			References:  []*Symbol{res},
		}
		a.file.Symbolcur.Symbols[res.Name] = res.SymbolCtx
	} else if st.Name != nil {
		err := &AnalysisErr{Errtype: SyntaxDataErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
	}
	//分析后面的数字表达式
	for _, exp := range []syntax.Expr{st.Init, st.Limit, st.Step} {
		switch res := a.analysisExpr(exp).(type) {
		case *Symbol: //名字
			if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
				res.SymbolCtx = sybif
				sybif.References = append(sybif.References, res)
				//检查类型是否为数字
				if !sybif.CurType.Contain(typeNumber) {
					err := &AnalysisErr{Errtype: LoopValNotNumErr}
					err.Scope = st.GetScope()
					err.insertInto(a)
				}
			} else {
				err := &AnalysisErr{Errtype: NoDefine}
				err.Scope = st.GetScope()
				err.insertInto(a)
			}
		case [][]TypeInfo: //检查类型是否为数字
			for _, info := range res[0] {
				if info == typeNumber {
					goto switchend
				}
			}
			err := &AnalysisErr{Errtype: LoopValNotNumErr}
			err.Scope = st.GetScope()
			err.insertInto(a)
		case TypeInfo: //检查类型是否为数字
			if res != typeNumber {
				err := &AnalysisErr{Errtype: LoopValNotNumErr}
				err.Scope = st.GetScope()
				err.insertInto(a)
			}
		case []*SymbolInfo: //检查类型是否为数字
			for _, info := range res {
				if info.CurType.Contain(typeNumber) {
					goto switchend
				}
			}
			err := &AnalysisErr{Errtype: LoopValNotNumErr}
			err.Scope = st.GetScope()
			err.insertInto(a)
		case nil:
		default:
			err := &AnalysisErr{Errtype: SyntaxDataErr}
			err.Scope = st.GetScope()
			err.insertInto(a)
		}
	}
switchend:
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
			res.SymbolCtx = &SymbolInfo{
				CurType:     NewTypeSetWithContent(typeNumber),
				Definitions: []*Symbol{res},
				References:  []*Symbol{res},
			}
			a.file.Symbolcur.Symbols[res.Name] = res.SymbolCtx
		} else if name != nil {
			err := &AnalysisErr{Errtype: SyntaxDataErr}
			err.Scope = st.GetScope()
			err.insertInto(a)
		}
	}
	//分析迭代表达式
	if len(st.Exprs) > 0 {
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
					err := &AnalysisErr{Errtype: LoopListParam1Err}
					err.Scope = st.Exprs[0].GetScope()
					err.insertInto(a)
				}
			} else {
				err := &AnalysisErr{Errtype: LoopListParam1Err}
				err.Scope = st.Exprs[0].GetScope()
				err.insertInto(a)
			}
		case *syntax.NameExpr:
			res := a.analysisNameExpr(expr)
			if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
				res.SymbolCtx = sybif
				sybif.References = append(sybif.References, res)
				//检查类型是否为函数
				count := 0
				for _, tp := range sybif.CurType.Types {
					if tp.TypeName() == "function" {
						count++
					}
				}
				if count == 0 {
					err := &AnalysisErr{Errtype: LoopListParam1Err}
					err.Scope = st.Exprs[0].GetScope()
					err.insertInto(a)
				}
			} else {
				err := &AnalysisErr{Errtype: NoDefine}
				err.Scope = st.Exprs[0].GetScope()
				err.insertInto(a)
			}
		}

		for i := 1; i < len(st.Exprs); i++ {
			if res, ok := a.analysisExpr(st.Exprs[i]).(*Symbol); ok {
				if sybif := a.file.Symbolcur.FindSymbol(res.Name); sybif != nil {
					res.SymbolCtx = sybif
					sybif.References = append(sybif.References, res)
				} else {
					err := &AnalysisErr{Errtype: NoDefine}
					err.Scope = st.Exprs[i].GetScope()
					err.insertInto(a)
				}
			}
		}
	} else {
		err := &AnalysisErr{Errtype: LoopListParamLoseErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
	}
	if len(st.Exprs) > 3 {
		err := &AnalysisErr{Errtype: LoopListParamMuchErr}
		err.Scope = st.GetScope()
		err.insertInto(a)
	}
	//分析内容
	for _, stmt := range st.Block {
		a.analysisStmt(stmt)
	}
}

func (a *Analysis) analysisFuncDefStmt(st *syntax.FuncDefStmt) {
	if st.Receiver != nil {
		switch expr := st.Receiver.(type) {
		case *syntax.GetItemExpr:
			if res := a.analysisGetItemExpr(expr); res != nil {

			}

		}
	} else {

	}
}
func (a *Analysis) analysisLocalFuncDefStmt(st *syntax.LocalFuncDefStmt) {
}
func (a *Analysis) analysisLocalVarDef(st *syntax.LocalVarDef) {
}
func (a *Analysis) analysisReturnStmt(st *syntax.ReturnStmt) {
}
func (a *Analysis) analysisErrorStmt(st *syntax.ErrorStmt) {
}
