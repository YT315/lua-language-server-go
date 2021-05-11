package syntax

type SyntaxErrBase string

const (
	StmtErr                 SyntaxErrBase = "语句错误"         //ErrorStmt
	LackRight               SyntaxErrBase = "缺少右值"         //FieldExpr	AssignStmt
	LackLeft                SyntaxErrBase = "缺少左值"         //AssignStmt
	LackIndex               SyntaxErrBase = "缺少索引"         //FieldExpr	GetItemExpr
	LackRightSquareBrackets SyntaxErrBase = "缺少右侧中括号"      //FieldExpr	GetItemExpr
	LackRightCurlyBrackets  SyntaxErrBase = "缺少右侧大括号"      //TableExpr
	Lackfuncbody            SyntaxErrBase = "缺少函数体"        //FuncDefExpr
	LackfuncName            SyntaxErrBase = "缺少函数名称"       //FuncCall	FuncDefStmt
	LackfuncArgs            SyntaxErrBase = "缺少参数"         //FuncCall
	LackField               SyntaxErrBase = "缺少字段"         //GetItemExpr
	LackObject              SyntaxErrBase = "缺少对象"         //GetItemExpr
	LabelIncomplete         SyntaxErrBase = "标签不完整"        //LabelStmt
	LackLabelName           SyntaxErrBase = "缺少标签名称"       //LabelStmt
	LackName                SyntaxErrBase = "缺少名称"         //LocalVarDef
	LackInitValue           SyntaxErrBase = "缺少初始值"        //LocalVarDef
	LackFunction            SyntaxErrBase = "缺少函数名称及内容"    //LocalFuncDefStmt
	LackFunctionContent     SyntaxErrBase = "缺少函数内容"       //LocalFuncDefStmt	FuncDefStmt
	LackFunctionkeyword     SyntaxErrBase = "缺少函数function" //LocalFuncDefStmt
	LackFunctionName        SyntaxErrBase = "缺少函数名称"       //LocalFuncDefStmt
	LackBlock               SyntaxErrBase = "缺少语句块"        //ForLoopListStmt
	LackEnd                 SyntaxErrBase = "缺少END"        //ForLoopListStmt	IfStmt	WhileStmt	DoEndStmt
	LackExpr                SyntaxErrBase = "缺少表达式"        //ForLoopListStmt
	LackForScope            SyntaxErrBase = "缺少循环范围"       //ForLoopNumStmt
	LackForStep             SyntaxErrBase = "缺少循环步进"       //ForLoopNumStmt
	LackForCond             SyntaxErrBase = "缺少循环条件"       //ForLoopNumStmt
	LackIfCond              SyntaxErrBase = "缺少判断条件"       //IfStmt
	LackUntilCond           SyntaxErrBase = "缺少until条件"    //RepeatStmt
	LackWhileCond           SyntaxErrBase = "缺少while条件"    //WhileStmt
	LackWhileBlock          SyntaxErrBase = "缺少while语句"    //WhileStmt
	LackGotoName            SyntaxErrBase = "缺少goto名称"     //GotoStmt
)
