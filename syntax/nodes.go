package syntax

//node
type (
	Node interface {
		nodeFlag()
	}

	nodeBase struct {
		Scope
		Err error
	}
)

func (n *nodeBase) nodeFlag() {}

//statements
type (
	Stmt interface {
		Node
		stmtFlag()
	}

	stmtBase struct {
		nodeBase
	}
)

func (*stmtBase) stmtFlag() {}

type (
	SemStmt struct {
		stmtBase
	}

	AssignStmt struct {
		stmtBase

		Left  []Expr
		Right []Expr
	}

	LabelStmt struct {
		stmtBase

		Name Expr
	}

	BreakStmt struct {
		stmtBase
	}

	GotoStmt struct {
		stmtBase

		Name Expr
	}

	DoEndStmt struct {
		stmtBase

		Block []Stmt
	}

	WhileStmt struct {
		stmtBase

		Condition Expr
		Block     []Stmt
	}

	RepeatStmt struct {
		stmtBase

		Condition Expr
		Block     []Stmt
	}

	IfStmt struct {
		stmtBase

		Condition Expr
		Then      []Stmt
		Else      []Stmt
	}

	ForLoopNumStmt struct {
		stmtBase

		Name  Expr
		Init  Expr
		Limit Expr
		Step  Expr
		Block []Stmt
	}

	ForLoopListStmt struct {
		stmtBase

		Names []Expr
		Exprs []Expr
		Block []Stmt
	}

	FuncDefStmt struct {
		stmtBase

		Name     Expr
		Receiver Expr

		Function Expr
	}

	LocalFuncDefStmt struct {
		stmtBase

		Name     Expr
		Function Expr
	}

	LocalVarDef struct {
		stmtBase

		Names []Expr
		Inits []Expr
	}

	ReturnStmt struct {
		stmtBase

		Exprs []Expr
	}

	ErrorStmt struct {
		stmtBase

		Info string
	}
)

//expressions
type (
	Expr interface {
		Node
		exprFlag()
	}

	exprBase struct {
		nodeBase
	}
)

func (*exprBase) exprFlag() {}

type (
	NameExpr struct {
		exprBase

		Value string
	}

	NilExpr struct {
		exprBase
	}

	FalseExpr struct {
		exprBase
	}

	TrueExpr struct {
		exprBase
	}

	NumberExpr struct {
		exprBase

		Value string
	}

	StringExpr struct {
		exprBase

		Value string
	}

	AnyExpr struct {
		exprBase
	}

	FuncDefExpr struct {
		exprBase

		Param Expr
		Block []Stmt
	}

	ParamExpr struct {
		exprBase

		Params []Expr
		IsAny  bool
	}

	GetItemExpr struct {
		exprBase

		Table Expr
		Key   Expr
	}

	TableExpr struct {
		exprBase

		Fields []Expr
	}

	FieldExpr struct {
		exprBase

		Key   Expr
		Value Expr
	}

	TwoOpExpr struct {
		exprBase

		Operator string
		Left     Expr
		Right    Expr
	}

	OneOpExpr struct {
		exprBase

		Operator string
		Target   Expr
	}
)

//spec
type (
	FuncCall struct {
		nodeBase

		Function Expr
		Receiver Expr
		Args     []Expr
	}
)

func (*FuncCall) stmtFlag() {}
func (*FuncCall) exprFlag() {}
