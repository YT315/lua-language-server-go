package analysis

import "lualsp/syntax"

//不会执行
func (a *Analysis) analysisSemStmt(st *syntax.SemStmt) {
}

func (a *Analysis) analysisAssignStmt(st *syntax.AssignStmt) {
	if len(st.Left) != len(st.Right) {
		//warrrrrrrrrrrrr
	}
	rlens := len(st.Right)
	for i, tvar := range st.Left {
		if rlens <= i { //判断此变量是否有右值
			continue
		}
		switch tp := tvar.(type) {
		case *syntax.NameExpr:
			name:=tp.Value 
			
			if 
			syb := &Symbol{
				Node: tvar,
				File: a.file,
			}

			//a.file.SymbolList.Symbols=append(a.file.SymbolList.Symbols,)
		}

	}

}
func (a *Analysis) analysisLabelStmt(st *syntax.LabelStmt) {
}
func (a *Analysis) analysisBreakStmt(st *syntax.BreakStmt) {
}
func (a *Analysis) analysisGotoStmt(st *syntax.GotoStmt) {
}
func (a *Analysis) analysisDoEndStmt(st *syntax.DoEndStmt) {
}
func (a *Analysis) analysisWhileStmt(st *syntax.WhileStmt) {
}
func (a *Analysis) analysisRepeatStmt(st *syntax.RepeatStmt) {
}
func (a *Analysis) analysisIfStmt(st *syntax.IfStmt) {
}
func (a *Analysis) analysisForLoopNumStmt(st *syntax.ForLoopNumStmt) {
}
func (a *Analysis) analysisForLoopListStmt(st *syntax.ForLoopListStmt) {
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
