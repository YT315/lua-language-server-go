%{
package syntax
%}

%union {
  token  Token

  stmts  []Stmt
  stmt   Stmt

  exprs  []Expr
  expr   Expr

  node   Node
}


//keyword
%token<token> TAnd TBreak TDo TElse TElseIf TEnd TFalse TFor TFunction TGoto TIf TIn TLocal TNil TNot TOr TRepeat TReturn TThen TTrue TUntil TWhile 

//符号         name   ==      ~=     <=       >=      <<     >>    //     ::     ..   ...   string  number
%token<token> TName TEqual TNequal TLequal TBequal TLmove TRmove TWdiv TTarget TConn TAny TString TNumber
%token<token> '+' '-' '*' '/' '%' '^' '#' '&' '~' '|' '<' '>' '=' '(' ')' '{' '}' '[' ']' ';' ':' ',' '.'

//Precedence
%left TOr
%left TAnd
%left '<' '>' TLequal TBequal TNequal TEqual
%left '|'
%left '~'
%left '&'
%left TLmove TRmove
%right TConn
%left '+' '-'
%left '*' '/' TWdiv '%'
%right UNARY /* not # - ~(unary) */
%right '^'

%type<stmts> chunk
%type<stmts> block
%type<stmts> blockaux
%type<stmt> stat
%type<stmt> returnstat
%type<exprs> varlist
%type<exprs> exprlist
%type<node> functioncall
%type<stmt> label
%type<stmts> elseifs
%type<stmt> funcname
%type<expr> funcnameaux
%type<expr> var
%type<exprs> namelist
%type<expr> expr
%type<expr> functiondef
%type<expr> funcbody
%type<expr> prefixexp
%type<exprs> args
%type<expr> parlist
%type<expr> tableconstructor
%type<exprs> fieldlist
%type<expr> field
%type<expr> fieldsep

%%

chunk: 
        block {
            $$ = $1
            if l, ok := lualex.(*Lexer); ok {
                l.Block = $$
            }
        }

block: 
        {
            $$ = []Stmt{}
        } |
        blockaux {
            $$ = $1
        } |
        returnstat {
            $$ = []Stmt{$1}
        } |
        blockaux returnstat {
            $$ = append($1, $2)
        }

blockaux:
        stat {
          if $1 != nil {
            $$ = []Stmt{$1}
          }else{
            $$ = []Stmt{}
          }
        } |
        blockaux stat {
          if $2 != nil{
            $$ = append($1, $2)
          }else{
            $$ = $1
          }
        }
stat:
        ';' {
            $$ = nil
        } |
        /*************** varlist ‘=’ explist *****************/
        varlist '=' exprlist {
            temp := &AssignStmt{Left: $1, Right: $3}
            temp.Start = $1[0].(*exprBase).Start
            temp.End = $3[len($3)-1].(*exprBase).End
            $$ = temp
        } |
        varlist '=' {
            temp := &AssignStmt{Left: $1, Right: nil}
            temp.Start=$1[0].(*exprBase).Start
            temp.End=$2.End
            temp.Err=&SyntaxErr{Info:"赋值表达式缺少右值"}
            temp.Err.Scope=$2.Scope
            $$ = temp
        } |  
        '=' exprlist {
            temp := &AssignStmt{Left: nil, Right: $2}
            temp.Start=$1.Start
            temp.End=$2[len($2)-1].(*exprBase).End
            temp.Err=&SyntaxErr{Info:"赋值表达式缺少左值"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        /*************** functioncall *****************/
        functioncall {
            $$ = $1
        } |
        /*************** label *****************/
        label {
            $$ = $1
        } |  
        /*************** TBreak *****************/
        TBreak {
            $$ = &BreakStmt{}
            $$.(*stmtBase).Scope = $1.Scope
        } |  
        /*************** goto Name *****************/
        TGoto TName {
            name := &NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp = &GotoStmt{Name:name}
            temp.Start=$1.Start
            temp.End = $2.End.
            $$ = temp
        } |
        TGoto {
            temp := &GotoStmt{Name:nil}
            temp.Scope = $1.Scope
            temp.Err=&SyntaxErr{Info:"缺少goto目标名称"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        /*************** do block end  *****************/
        TDo block TEnd {
            temp := &DoEndStmt{Block: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            $$ = temp
        } |
        TDo block {
            temp := &DoEndStmt{Block: $2}
            temp.Start=$1.Start
            if len($2)>0 {
                temp.End = $2[len($2)-1].(*stmtBase).End
            }else{
                temp.End = $1.End
            }
            temp.Err=&SyntaxErr{Info:"缺少End"}
            temp.Err.Start = temp.End
            temp.Err.End = temp.End
            $$ = temp
        } |
        /*************** while exp do block end  *****************/ 
        TWhile expr TDo block TEnd {
            temp := &WhileStmt{Condition: $2, Block: $4}
            temp.Start=$1.Start
            temp.End = $5.End
            $$ = temp
        } |
        TWhile expr TDo block {
            temp := &WhileStmt{Condition: $2, Block: $4}
            temp.Start=$1.Start
            if len($4)>0 {
                temp.End = $4[len($4)-1].(*stmtBase).End
            }else{
                temp.End = $3.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TWhile TDo block TEnd {
            temp := &WhileStmt{Condition: nil, Block: $3}
            temp.Start=$1.Start
            temp.End = $4.End
            temp.Err=&SyntaxErr{Info:"缺少条件表达式"}
            temp.Err.Start=$1.End
            temp.Err.End=$2.Start     
            $$ = temp
        } |
        TWhile expr {
            temp := &WhileStmt{Condition: $2, Block: nil}
            temp.Start=$1.Start
            temp.End = $2.(*exprBase).End
            temp.Err=&SyntaxErr{Info:"缺少语句do..end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TWhile {
            temp := &WhileStmt{Condition: nil, Block: nil}
            temp.Scope=$1.Scope
            temp.Err=&SyntaxErr{Info:"缺少语句expr do..end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        /*************** repeat block until exp *****************/ 
        TRepeat block TUntil expr {
            temp := &RepeatStmt{Condition: $4, Block: $2}
            temp.Start=$1.Start
            temp.End = $4.(*exprBase).End
            $$ = temp
        } |
        TRepeat block TUntil {
            temp := &RepeatStmt{Condition: nil, Block: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            temp.Err=&SyntaxErr{Info:"缺少语句条件表达式"}
            temp.Err.Scope=$3.Scope
            $$ = temp
        } |
        TRepeat block {
            temp := &RepeatStmt{Condition: nil, Block: $2}
            temp.Start=$1.Start
            if len($2)>0 {
                temp.End = $2[len($2)-1].(*stmtBase).End
            }else{
                temp.End = $1.End
            }
            temp.Err=&SyntaxErr{Info:"缺少条件以及Until"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        /*************** TIf expr TThen block elseifs TEnd *****************/ 
        TIf expr TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.(*stmtBase).Start=$1.Start
            $$.(*stmtBase).End = $6.End
        } |
        TIf expr TThen block elseifs {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            if len($5)>0{
                temp.End = $5[len($5)-1].(*stmtBase).End
            }else if len($4)>0 {
                temp.End = $4[len($4)-1].(*stmtBase).End
            }else{
                temp.End = $3.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Scope=$1.Scope
        } |
        TIf TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            temp.End = $5.End
            temp.Err=&SyntaxErr{Info:"缺少条件表达式"}
            temp.Err.Start=$1.End
            temp.Err.End=$2.Start   
        } |
        TIf block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            temp.End = $4.End
            temp.Err=&SyntaxErr{Info:"缺少条件表达式及then"}
            temp.Err.Scope=$1.Scope
        } |
        //err
        TIf block elseifs {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            if len($3)>0{
                temp.End = $3[len($3)-1].(*stmtBase).End
            }else if len($2)>0{
                temp.End = $2[len($2)-1].(*stmtBase).End
            }else{
                temp.End=$1.End
            }
            temp.Err=&SyntaxErr{Info:"缺少条件,then和end"}
            temp.Err.Scope=temp.Scope
        } |
        /*************** TIf expr TThen block elseifs TElse block TEnd *****************/ 
        TIf expr TThen block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $7
            $$.(*stmtBase).Start = $1.Start
            $$.(*stmtBase).End = $8.End
        } |
        TIf TThen block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $6
            temp :=$$.(*stmtBase)
            temp.Start = $1.Start
            temp.End = $7.End
            temp.Err=&SyntaxErr{Info:"缺少条件表达式"}
            temp.Err.Scope=temp.Scope
        } |
        TIf expr TThen block elseifs TElse block {
            $$ = &IfStmt{Condition: $2, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $6
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            if len($7)>0 {
                temp.End = $7[len($7)-1].(*stmtBase).End
            }else{
                temp.End = $6.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Start=temp.End
            temp.Err.End=temp.End
        } |
        TIf block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $5
            temp :=$$.(*stmtBase)
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Info:"缺少条件表达式及then"}
            temp.Err.Scope=$1.Scope
        } |
        /*************** TFor TName '=' expr ',' expr TDo block TEnd *****************/
        TFor TName '=' expr ',' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Block: $8}
            $$.(*stmtBase).Start=$1.Start
            $$.(*stmtBase).End = $9.End
        } |
        TFor TName '=' expr ',' expr TDo block {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Block: $8}
            temp.Start=$1.Start
            if len($8)>0 {
                temp.End = $8[len($8)-1].End
            }else{
                temp.End = $7.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TFor '=' expr ',' expr TDo block TEnd {
            temp := &ForLoopNumStmt{Name: nil, Init: $3, Limit: $5, Block: $7}
            temp.Start=$1.Start
            temp.End = $8.End
            temp.Err=&SyntaxErr{Info:"缺少名称"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TFor TDo block TEnd {
            temp := &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            temp.Start=$1.Start
            temp.End = $4.End
            temp.Err=&SyntaxErr{Info:"缺少循环条件"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TFor TDo block {
            temp := &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            temp.Start=$1.Start
            if len($3)>0 {
                temp.End = $3[len($3)-1].End
            }else{
                temp.End = $2.End
            }
            temp.Err=&SyntaxErr{Info:"缺少循环条件及end"}
            temp.Err.Scope=temp.Scope
            $$ = temp
        } |
        TFor TName TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: nil, Limit: nil, Block: $4}
            temp.Start=$1.Start
            temp.End = $5.End
            temp.Err=&SyntaxErr{Info:"缺少循环范围"}
            temp.Err.Scope=$2.Scope
            $$ = temp
        } |
        TFor TName '=' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: nil, Limit: nil, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Info:"缺少循环范围"}
            temp.Err.Scope=$2.Scope
            $$ = temp
        } |
        TFor TName '=' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: nil, Block: $6}
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Info:"缺少循环终点"}
            temp.Err.Scope=$4.Scope
            $$ = temp
        } |
        TFor TName '=' expr ',' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: nil, Block: $7}
            temp.Start=$1.Start
            temp.End = $7.End
            temp.Err=&SyntaxErr{Info:"缺少循环终点"}
            temp.Err.Scope=$5.Scope
            $$ = temp
        } |
        /*************** TFor TName '=' expr ',' expr ',' expr TDo block TEnd *****************/
        TFor TName '=' expr ',' expr ',' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:$8, Block: $10}
            temp.Start=$1.Start
            temp.End = $11.End
            $$ = temp
        } |
        TFor TName '=' expr ',' expr ',' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:nil, Block: $9}
            temp.Start=$1.Start
            temp.End = $10.End
            temp.Err=&SyntaxErr{Info:"缺少步进值"}
            temp.Err.Start=$7.End
            temp.Err.End=$8.Start
            $$ = temp
        } |
        TFor TName '=' expr ',' expr ',' expr TDo block {
            name:=&NameExpr{Value:$2.Str}
            name.Scope = $2.Scope
            temp := &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:$8, Block: $10}
            temp.Start=$1.Start
            if len($10)>0{
                temp.End = $10[len($10)-1].(*stmtBase).End
            }else{
                temp.End = $9.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        /*************** TFor namelist TIn exprlist TDo block TEnd *****************/
        TFor namelist TIn exprlist TDo block TEnd {
            $$ = &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            $$.(*stmtBase).Start=$1.Start
            $$.(*stmtBase).End = $7.End
        } |
        TFor TIn exprlist TDo block TEnd {
            temp := &ForLoopListStmt{Exprs:$3, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Info:"缺少迭代表达式"}
            temp.Err.Scope=$2.Scope
            $$ = temp
        } |
        TFor namelist TIn TDo block TEnd {
            temp := &ForLoopListStmt{Names:$2, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Info:"缺少迭代对象表达式"}
            temp.Err.Scope=$2.Scope
            $$ = temp
        } |
        TFor namelist TIn exprlist TDo block {
            temp := &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            temp.Start=$1.Start
            if len($6)>0{
                temp.End = $6[len($6)-1].End
            }else{
                temp.End = $5.End
            }
            temp.Err=&SyntaxErr{Info:"缺少end"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |
        TFor namelist TIn exprlist{
            temp := &ForLoopListStmt{Names:$2, Exprs:$4}
            temp.Start=$1.Start
            if len($4)>0{
                temp.End = $4[len($4)-1].End
            }else{
                temp.End = $3.End
            }
            temp.Err=&SyntaxErr{Info:"缺少执行语句end"}
            temp.Err.Scope=$1.Scope
        } |
        /*************** TFunction funcname funcbody *****************/
        TFunction funcname funcbody {
            $$ = $2
            $$.(*FuncDefStmt).Function= $3
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        TFunction funcname {
            $$ = $2
            $$.Start=$1.Start
            $$.Err=&SyntaxErr{Info:"缺少函数体"}
            $$.Err.Scope=$1.Scope
        } |
        TFunction funcbody {
            $$ = &FuncDefStmt{}
            $$.(*FuncDefStmt).Function= $2
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少函数名"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** TLocal TFunction TName funcbody *****************/
        TLocal TFunction TName funcbody {
            name:=&NameExpr{Value:$3.Str}
            $$ = &LocalFuncDefStmt{Name: name, Function: $4}
            $$.Start=$1.Start
            $$.End = $4.End
        } | 
        TLocal TFunction funcbody {
            $$ = &LocalFuncDefStmt{Function: $3}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少函数名"}
            $$.Err.Scope=$2.Scope
        } | 
        TLocal TName funcbody {
            name:=&NameExpr{Value:$2.Str}
            $$ = &LocalFuncDefStmt{Name: name, Function: $3}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少function"}
            $$.Err.Scope=$1.Scope
        } | 
        TLocal TFunction TName {
            name:=&NameExpr{Value:$3.Str}
            $$ = &LocalFuncDefStmt{Name: name}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少函数体"}
            $$.Err.Scope=$3.Scope
        } | 
        TLocal TFunction {
            $$ = &LocalFuncDefStmt{}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少函数内容"}
            $$.Err.Scope=$2.Scope
        } | 
        /*************** TLocal namelist '=' exprlist *****************/
        TLocal namelist '=' exprlist {
            $$ = &LocalVarDef{Names: $2, Inits:$4}
            $$.Start=$1.Start
            $$.End = $4.End
        } |
        TLocal namelist '=' {
            $$ = &LocalVarDef{Names: $2}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少初始值"}
            $$.Err.Scope=$3.Scope
        } |
        TLocal '=' exprlist {
            $$ = &LocalVarDef{Inits:$3}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少变量名称"}
            $$.Err.Scope=$2.Scope
        } |
        /*************** TLocal namelist *****************/
        TLocal namelist {
            $$ = &LocalVarDef{Names: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        }|
        TLocal {
            $$ = &LocalVarDef{}
            $$.Scope=$1.Scope
            $$.Err=&SyntaxErr{Info:"缺少变量名称"}
            $$.Err.Scope=$1.Scope
        }|
        error{
            $$ = &ErrorStmt{Info:"解析错误"}
            tk:=lualex.(Lexer).Token
            $$.Err=&SyntaxErr{Info:tk.Str+"附近解析错误"}
            $$.Err.Scope=tk.Scope
        }

elseifs: 
        {
            $$ = []Stmt{}
        } | 
        elseifs TElseIf expr TThen block {
            ifstmt:=&IfStmt{Condition: $3, Then: $5}
            ifstmt.Start=$2.Start
            ifstmt.End = $5.End
            $$ = append($1, ifstmt)
        }
        
returnstat:
        TReturn {
            $$ = &ReturnStmt{Exprs:nil}
            $$.Scope=$1.Scope
        } |
        TReturn exprlist {
            $$ = &ReturnStmt{Exprs:$2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        TReturn exprlist ';' {
            $$ = &ReturnStmt{Exprs:$2}
            $$.Start=$1.Start
            $$.End = $3.End
        }

label:  
        TTarget TName TTarget {
            name := &NameExpr{Value:$2.Str}
            name.Scope:$2.Scope
            $$ = &LabelStmt{Name: name}
            $$.Start=$1.Start
            $$.End=$3.End
        }|
        TName TTarget {
            name := &NameExpr{Value:$2.Str}
            name.Scope:$1.Scope
            $$ = &LabelStmt{Name: name}
            $$.Start=$1.Start
            $$.End=$2.End
            $$.Err=&SyntaxErr{Info:"标签缺少左侧符号"}
            $$.Err.Scope=$1.Scope
        }|
        TTarget TName {
            name:=&NameExpr{Value:$2.Str}
            name.Scope:$2.Scope
            $$ = &LabelStmt{Name: name}
            $$.Start=$1.Start
            $$.End=$2.End
            $$.Err=&SyntaxErr{Info:"标签缺少右侧符号"}
            $$.Err.Scope=$2.Scope
        }|
        TTarget TTarget {
            $$ = &LabelStmt{Name: nil}
            $$.Start=$1.Start
            $$.End=$2.End
            $$.Err=&SyntaxErr{Info:"标签缺少名称"}
            $$.Err.Scope=$$.Scope
        }|
        TTarget {
            $$ = &LabelStmt{Name: nil}
            $$.Scope=$1.Scope
            $$.Err=&SyntaxErr{Info:"标签不完整"}
            $$.Err.Scope=$$.Scope
        }
        
funcname: 
        funcnameaux {
            $$ = &FuncDefStmt{Name: $1, Receiver: nil}
            $$.Scope=$1.Scope
        } |
        funcnameaux ':' TName {
            name:= &NameExpr{Value:$3.Str}
            name.Scope=$3.Scope
            $$ = &FuncDefStmt{Name: name, Receiver:$1}
            $$.Start=$1.Start
            $$.End = $3.End
        }|
        funcnameaux ':' {
            $$ = &FuncDefStmt{Name: name}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少接受者"}
            $$.Err.Scope=$2.Scope
        }

funcnameaux:
        TName {
            $$ = &NameExpr{Value:$1.Str}
            $$.Scope=$1.Scope
        } | 
        funcnameaux '.' TName {
            name:= &NameExpr{Value:$3.Str}
            name.Scope=$1.Scope
            $$ = &GetItemExpr{Table:$1, Key:name}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        funcnameaux '.' {
            $$ = &GetItemExpr{Table:$1}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少子项目"}
            $$.Err.Scope=$2.Scope
        }

varlist:
        var {
            $$ = []Expr{$1}
        } | 
        varlist ',' var {
            $$ = append($1, $3)
        } 

var:
        TName {
            $$ = &NameExpr{Value:$1.Str}
            $$.Scope=$1.Scope
        } |
        prefixexp '[' expr ']' {
            $$ = &GetItemExpr{Table:$1, Key:$3}
            $$.Start=$1.Start
            $$.End = $4.End
        } | 
        prefixexp '[' ']' {
            $$ = &GetItemExpr{Table:$1}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少索引"}
            $$.Err.Scope=$2.Scope
        } |
        prefixexp '[' expr {
            $$ = &GetItemExpr{Table:$1, Key:$3}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少右括号"}
            $$.Err.Scope=$3.Scope
        } |
        prefixexp '[' {
            $$ = &GetItemExpr{Table:$1}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少索引"}
            $$.Err.Scope=$2.Scope
        } |
        prefixexp '.' TName {
            name:= &NameExpr{Value:$3.Str}
            $$.Scope=$3.Scope
            $$ = &GetItemExpr{Table:$1, Key:name}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        prefixexp '.' {
            $$ = &GetItemExpr{Table:$1}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少子项目"}
            $$.Err.Scope=$2.Scope
        }

namelist:
        TName {
            name:= &NameExpr{Value:$1.Str}
            $$.Scope=$1.Scope
            $$ = []Expr{name}
        } | 
        namelist ','  TName {
            name:= &NameExpr{Value:$3.Str}
            $$.Scope=$3.Scope
            $$ = append($1, name)
        }

exprlist:
        expr {
            $$ = []Expr{$1}
        } |
        exprlist ',' expr {
            $$ = append($1, $3)
        }

expr:
        TNil {
            $$ = &NilExpr{}
            $$.Scope=$1.Scope
        } | 
        TFalse {
            $$ = &FalseExpr{}
            $$.Scope=$1.Scope
        } | 
        TTrue {
            $$ = &TrueExpr{}
            $$.Scope=$1.Scope
        } | 
        TNumber {
            $$ = &NumberExpr{Value: $1.Str}
            $$.Scope=$1.Scope
        } | 
        TString {
            $$ = &StringExpr{Value: $1.Str}
            $$.Scope=$1.Scope
        } |
        TAny {
            $$ = &AnyExpr{}
            $$.Scope=$1.Scope
        } |
        functiondef {
            $$ = $1
        } |
        prefixexp {
            $$ = $1
        } |
        tableconstructor {
            $$ = $1
        } |
        expr '+' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '-' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '*' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '/' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TWdiv expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '^' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '%' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '&' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '~' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '|' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TRmove expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TLmove expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TConn expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '<' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TLequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr '>' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TBequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TEqual expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TNequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TAnd expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        expr TOr expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        '-' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        TNot expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        '#' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        '~' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        '(' expr ')' {
            $$ = $2
            $$.Start=$1.Start
            $$.End = $3.End
        }

prefixexp:
        var {
            $$ = $1
        } |
        functioncall {
            $$ = $1
        } 

functioncall:
        prefixexp args {
            $$ = &FuncCall{Function: $1, Args: $2}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        prefixexp ':' TName args {
            name:= &NameExpr{Value: $3.Str}
            name.Scope=$3.Scope
            $$ = &FuncCall{Function: name, Receiver: $1, Args: $4}
            $$.Start=$1.Start
            $$.End = $4.End
        } | 
        prefixexp ':' TName  {
            name:= &NameExpr{Value: $3.Str}
            name.Scope=$3.Scope
            $$ = &FuncCall{Function: name, Receiver: $1,}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少函数调用参数(args)"}
            $$.Err.Scope=$2.Scope
        } | 
        prefixexp ':'  {
            $$ = &FuncCall{ Receiver: $1}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少函数调用"}
            $$.Err.Scope=$2.Scope
        }

args:
        '(' ')' {
            $$ = []Expr{}
        } |
        '(' exprlist ')' {
            $$ = $2
        } |
        tableconstructor {
            $$ = []Expr{$1}
        } | 
        TString {
            str := &StringExpr{Value: $1.Str}
            str.Scope =$1.Scope
            $$ = []Expr{str}
        }

functiondef:
        TFunction funcbody {
            $$ = $2
        } |
        TFunction {
            $$ =  &FuncDefExpr{}
            $$.Scope =$1.Scope
            $$.Err=&SyntaxErr{Info:"未定义函数体"}
            $$.Err.Scope=$1.Scope
        }

funcbody:
        '(' parlist ')' block TEnd {
            $$ = &FuncDefExpr{Param: $2, Block: $4}
            $$.Start=$1.Start
            $$.End = $5.End
        } | 
        '(' ')' block TEnd {
            $$ = &FuncDefExpr{Param: nil, Block: $3}
            $$.Start=$1.Start
            $$.End = $4.End
        }

parlist:
        TAny {
            $$ = &ParamExpr{IsAny: true}
            $$.Scope =$1.Scope
        } | 
        namelist {
            $$ = &ParamExpr{Params: $1, IsAny: false}
            $$.Scope =$1.Scope
        } | 
        namelist ',' TAny {
            $$ = &ParamExpr{Params: $1, IsAny: true}
            $$.Start=$1.Start
            $$.End = $3.End
        }


tableconstructor:
        '{' '}' {
            $$ = &TableExpr{Fields: []Expr{}}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        '{' fieldlist '}' {
            $$ = &TableExpr{Fields: $2}
            $$.Start=$1.Start
            $$.End = $3.End
        }


fieldlist:
        field {
            $$ = []Expr{$1}
        } | 
        fieldlist fieldsep field {
            $$ = append($1, $3)
        } | 
        fieldlist fieldsep {
            $$ = $1
        }

field:
        TName '=' expr {
            name:= &NameExpr{Value: $1.Str}
            name.Scope =$1.Scope
            $$ = &FieldExpr{Key: name, Value: $3}
            $$.Start=$1.Start
            $$.End = $3.End
        } | 
        '[' expr ']' '=' expr {
            $$ = &FieldExpr{Key: $2, Value: $5}
            $$.Start=$1.Start
            $$.End = $5.End
        } |
        expr {
            $$ = &FieldExpr{Value: $1}
            $$.Scope =$1.Scope
        }

fieldsep:
        ',' {
            $$ = &NameExpr{Value: $1.Str}
        } | 
        ';' {
            $$ = &NameExpr{Value: $1.Str}
        }
%%