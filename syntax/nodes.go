package syntax

import (
	"fmt"
	"reflect"
	"strings"
)

//node
type (
	Node interface {
		scope() Scope
		start() Pos
		end() Pos
	}

	nodeBase struct {
		Scope
		Err *SyntaxErr
	}
)

func (n *nodeBase) scope() Scope { return n.Scope }
func (n *nodeBase) start() Pos   { return n.Start }
func (n *nodeBase) end() Pos     { return n.End }

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
		setBracket(bool)
	}

	exprBase struct {
		nodeBase
		HaveBracket bool
	}
)

func (eb *exprBase) setBracket(b bool) {
	eb.HaveBracket = b
}

type (
	NameExpr struct {
		exprBase
		Type  Expr
		Value string
	}

	STypeExpr struct {
		exprBase
		Value string
	}

	ATypeExpr struct {
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

		Value float64
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

		Param  Expr
		Result Expr
		Block  []Stmt
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

func (*FuncCall) stmtFlag()       {}
func (*FuncCall) setBracket(bool) {}

//Traversal 遍历抽象语法树，绘制出树结构
func Traversal(node interface{}, vist func(Node)) {
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Invalid {
		fmt.Printf("nill\n")
		return
	}
	if n, ok := node.(Node); ok {
		vist(n)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t := v.Type()
		//fmt.Printf("\n")
		for i := 0; i < t.NumField(); i++ {
			//基类不进
			nm := t.Field(i).Name
			if strings.Index(nm, "Base") != -1 {
				continue
			}
			sv := v.Field(i)
			switch sv.Kind() {
			case reflect.Slice:
				for j := 0; j < sv.Len(); j++ {
					ifc := sv.Index(j).Interface()
					if _, ok := ifc.(Node); ok {
						Traversal(ifc, vist)
					}
				}
			case reflect.Ptr:
				ifc := sv.Interface()
				if _, ok := ifc.(Node); ok {
					Traversal(ifc, vist)
				}
			case reflect.Interface:
				ifc := sv.Interface()
				if _, ok := ifc.(Node); ok {
					Traversal(ifc, vist)
				}
			default:
			}
		}
	} else {
		panic(v.Kind())
	}

}
