package analysis

import "lualsp/syntax"

//Analysis 语义分析器
type Analysis struct {
	previous *Analysis //对象包含依赖,分析过程通过单项链表连接
	file     *File     //正在分析的文件指针
}

func (a *Analysis) star() {
	//遍历ast
	for _, stmt := range a.file.Ast {
		a.analysisStmt(stmt)
	}
}

//分析语句
func (a *Analysis) analysisStmt(st syntax.Stmt) {
	switch tp := st.(type) {
	case *syntax.AssignStmt:
		a.analysisAssignStmt(tp)
	case *syntax.LabelStmt:
		a.analysisLabelStmt(tp)
	case *syntax.BreakStmt:
		a.analysisBreakStmt(tp)
	case *syntax.GotoStmt:
		a.analysisGotoStmt(tp)
	case *syntax.DoEndStmt:
		a.analysisDoEndStmt(tp)
	case *syntax.WhileStmt:
		a.analysisWhileStmt(tp)
	case *syntax.RepeatStmt:
		a.analysisRepeatStmt(tp)
	case *syntax.IfStmt:
		a.analysisIfStmt(tp)
	case *syntax.ForLoopNumStmt:
		a.analysisForLoopNumStmt(tp)
	case *syntax.ForLoopListStmt:
		a.analysisForLoopListStmt(tp)
	case *syntax.FuncDefStmt:
		a.analysisFuncDefStmt(tp)
	case *syntax.LocalFuncDefStmt:
		a.analysisLocalFuncDefStmt(tp)
	case *syntax.LocalVarDef:
		a.analysisLocalVarDef(tp)
	case *syntax.ReturnStmt:
		a.analysisReturnStmt(tp)
	case *syntax.ErrorStmt:
		a.analysisErrorStmt(tp)
	case *syntax.FuncCall:
		a.analysisFuncCall(tp)
	}
}

//分析表达式
func (a *Analysis) analysisExpr(ep syntax.Expr) interface{} {
	switch tp := ep.(type) {
	//return TypeInfo
	case *syntax.NilExpr:
		return a.analysisNilExpr(tp)
	case *syntax.FalseExpr:
		return a.analysisFalseExpr(tp)
	case *syntax.TrueExpr:
		return a.analysisTrueExpr(tp)
	case *syntax.NumberExpr:
		return a.analysisNumberExpr(tp)
	case *syntax.StringExpr:
		return a.analysisStringExpr(tp)
	case *syntax.AnyExpr:
		return a.analysisAnyExpr(tp)
	case *syntax.FuncDefExpr:
		return a.analysisFuncDefExpr(tp)
	case *syntax.TableExpr:
		return a.analysisTableExpr(tp)
	case *syntax.TwoOpExpr:
		return a.analysisTwoOpExpr(tp)
	case *syntax.OneOpExpr:
		return a.analysisOneOpExpr(tp)
	//return [][]TypeInfo
	case *syntax.FuncCall:
		return a.analysisFuncCall(tp)
	//return *Symbol
	case *syntax.NameExpr:
		return a.analysisNameExpr(tp)
	//result []*SymbolInfo
	case *syntax.GetItemExpr:
		return a.analysisGetItemExpr(tp)
	default:
		return nil
	}
}
func (a *Analysis) analysisFuncCall(n *syntax.FuncCall) [][]TypeInfo { //第一维索引,第二维类型
	return nil
}
