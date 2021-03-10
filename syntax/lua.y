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
            $$ = &AssignStmt{Left: $1, Right: $3}
            $$.Start=$1[0].Start
            $$.End=$3[len($3)-1].End
        } |
        varlist '=' {
            $$ = &AssignStmt{Left: $1, Right: nil}
            $$.Start=$1[0].Start
            $$.End=$2.End
            $$.Err=&SyntaxErr{Info:"赋值表达式缺少右值"}
            $$.Err.Scope=$2.Scope
        } |  
        '=' exprlist {
            $$ = &AssignStmt{Left: nil, Right: $2}
            $$.Start=$1.Start
            $$.End=$2[len($2)-1].End
            $$.Err=&SyntaxErr{Info:"赋值表达式缺少左值"}
            $$.Err.Scope=$1.Scope
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
            $$.Scope = $1.Scope
        } |  
        /*************** goto Name *****************/
        TGoto TName {
            $$ = &GotoStmt{Name:&NameExpr{Value:$2.Str}}
            $$.Start=$1.Start
            $$.End = $2.End
        } |
        TGoto {
            $$ = &GotoStmt{Name:nil}
            $$.Scope = $1.Scope
            $$.Err=&SyntaxErr{Info:"缺少goto目标名称"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** do block end  *****************/
        TDo block TEnd {
            $$ = &DoEndStmt{Block: $2}
            $$.Start=$1.Start
            $$.End = $3.End
        } |
        TDo block {
            $$ = &DoEndStmt{Block: $2}
            $$.Start=$1.Start
            $$.End = $2[len($2)-1].End
            $$.Err=&SyntaxErr{Info:"缺少End"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** while exp do block end  *****************/ 
        TWhile expr TDo block TEnd {
            $$ = &WhileStmt{Condition: $2, Block: $4}
            $$.Start=$1.Start
            $$.End = $5.End
        } |
        TWhile expr TDo block {
            $$ = &WhileStmt{Condition: $2, Block: $4}
            $$.Start=$1.Start
            $$.End = $4[len($4)-1].End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Scope=$1.Scope
        } |
        TWhile TDo block TEnd {
            $$ = &WhileStmt{Condition: nil, Block: $3}
            $$.Start=$1.Start
            $$.End = $4.End
            $$.Err=&SyntaxErr{Info:"缺少条件表达式"}
            $$.Err.Start=$1.End
            $$.Err.End=$2.Start     
        } |
        TWhile expr {
            $$ = &WhileStmt{Condition: $2, Block: nil}
            $$.Start=$1.Start
            $$.End = $2.End
            $$.Err=&SyntaxErr{Info:"缺少语句do..end"}
            $$.Err.Scope=$1.Scope
        } |
        TWhile {
            $$ = &WhileStmt{Condition: nil, Block: nil}
            $$.Scope=$1.Scope
            $$.Err=&SyntaxErr{Info:"缺少语句expr do..end"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** repeat block until exp *****************/ 
        TRepeat block TUntil expr {
            $$ = &RepeatStmt{Condition: $4, Block: $2}
            $$.Start=$1.Start
            $$.End = $4.End
        } |
        TRepeat block TUntil {
            $$ = &RepeatStmt{Condition: nil, Block: $2}
            $$.Start=$1.Start
            $$.End = $3.End
            $$.Err=&SyntaxErr{Info:"缺少语句条件表达式"}
            $$.Err.Scope=$3.Scope
        } |
        TRepeat block {
            $$ = &RepeatStmt{Condition: nil, Block: $2}
            $$.Start=$1.Start
            $$.End = $2[len($2)-1].End
            $$.Err=&SyntaxErr{Info:"缺少条件以及Until"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** TIf expr TThen block elseifs TEnd *****************/ 
        TIf expr TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.Start=$1.Start
            $$.End = $6.End
        } |
        //err
        TIf expr TThen block elseifs {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.Start=$1.Start
            $$.End = $5[len($5)-1].End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Scope=$1.Scope
        } |
        //err
        TIf TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.Start=$1.Start
            $$.End = $5.End
            $$.Err=&SyntaxErr{Info:"缺少条件表达式"}
            $$.Err.Start=$1.End
            $$.Err.End=$2.Start   
        } |
        //err
        TIf block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.Start=$1.Start
            $$.End = $4.End
            $$.Err=&SyntaxErr{Info:"缺少条件表达式及then"}
            $$.Err.Scope=$1.Scope
        } |
        //err
        TIf block elseifs {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.Start=$1.Start
            if len($3)>0{
                $$.End = $3[len($3)-1].End
            }else if len($2)>0{
                $$.End = $2[len($2)-1].End
            }else{
                $$.End=$1.End
            }
            $$.Err=&SyntaxErr{Info:"缺少条件,then和end"}
            $$.Err.Scope=$$.Scope
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
            $$.Start=$1.Start
            $$.End = $8.End
        } |
        TIf TThen block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $6
            $$.Start=$1.Start
            $$.End = $7.End
            $$.Err=&SyntaxErr{Info:"缺少条件表达式"}
            $$.Err.Scope=$$.Scope
        } |
        TIf expr TThen block elseifs TElse block {
            $$ = &IfStmt{Condition: $2, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $6
            $$.Start=$1.Start
            $$.End = $6[len($6)-1].End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Start=$$.End
            $$.Err.End=$$.End
        } |
        TIf block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $5
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少条件表达式及then"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** TFor TName '=' expr ',' expr TDo block TEnd *****************/
        TFor TName '=' expr ',' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Block: $8}
            $$.Start=$1.Start
            $$.End = $9.End
        } |
        TFor TName '=' expr ',' expr TDo block {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Block: $8}
            $$.Start=$1.Start
            $$.End = $8[len($8)-1].End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Scope=$1.Scope
        } |
        TFor '=' expr ',' expr TDo block TEnd {
            $$ = &ForLoopNumStmt{Name: nil, Init: $3, Limit: $5, Block: $7}
            $$.Start=$1.Start
            $$.End = $8.End
            $$.Err=&SyntaxErr{Info:"缺少名称"}
            $$.Err.Scope=$1.Scope
        } |
        TFor TDo block TEnd {
            $$ = &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            $$.Start=$1.Start
            $$.End = $4.End
            $$.Err=&SyntaxErr{Info:"缺少循环条件"}
            $$.Err.Scope=$1.Scope
        } |
        TFor TDo block {
            $$ = &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            $$.Start=$1.Start
            $$.End = $3[len($3)-1].End
            $$.Err=&SyntaxErr{Info:"缺少循环条件及end"}
            $$.Err.Scope=$$.Scope
        } |
        TFor TName TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: nil, Limit: nil, Block: $4}
            $$.Start=$1.Start
            $$.End = $5.End
            $$.Err=&SyntaxErr{Info:"缺少循环范围"}
            $$.Err.Scope=$2.Scope
        } |
        TFor TName '=' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: nil, Limit: nil, Block: $5}
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少循环范围"}
            $$.Err.Scope=$2.Scope
        } |
        TFor TName '=' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: nil, Block: $6}
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少循环终点"}
            $$.Err.Scope=$4.Scope
        } |
        TFor TName '=' expr ',' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: nil, Block: $7}
            $$.Start=$1.Start
            $$.End = $7.End
            $$.Err=&SyntaxErr{Info:"缺少循环终点"}
            $$.Err.Scope=$5.Scope
        } |
        /*************** TFor TName '=' expr ',' expr ',' expr TDo block TEnd *****************/
        TFor TName '=' expr ',' expr ',' expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:$8, Block: $10}
            $$.Start=$1.Start
            $$.End = $11.End
        } |
        TFor TName '=' expr ',' expr ',' TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:nil, Block: $9}
            $$.Start=$1.Start
            $$.End = $10.End
            $$.Err=&SyntaxErr{Info:"缺少步进值"}
            $$.Err.Start=$7.End
            $$.Err.End=$8.Start
        } |
        TFor TName '=' expr ',' expr ',' expr TDo block {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $6, Step:$8, Block: $10}
            $$.Start=$1.Start
            $$.End = $10[len($10)-1].End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Scope=$1.Scope
        } |
        /*************** TFor namelist TIn exprlist TDo block TEnd *****************/
        TFor namelist TIn exprlist TDo block TEnd {
            $$ = &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            $$.Start=$1.Start
            $$.End = $7.End
        } |
        TFor TIn exprlist TDo block TEnd {
            $$ = &ForLoopListStmt{Exprs:$3, Block: $5}
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少迭代表达式"}
            $$.Err.Scope=$2.Scope
        } |
        TFor namelist TIn TDo block TEnd {
            $$ = &ForLoopListStmt{Names:$2, Block: $5}
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少迭代对象表达式"}
            $$.Err.Scope=$2.Scope
        } |
        TFor namelist TIn exprlist TDo block {
            $$ = &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            $$.Start=$1.Start
            $$.End = $6.End
            $$.Err=&SyntaxErr{Info:"缺少end"}
            $$.Err.Scope=$1.Scope
        } |
        TFor namelist TIn exprlist{
            $$ = &ForLoopListStmt{Names:$2, Exprs:$4}
            $$.Start=$1.Start
            $$.End = $4.End
            $$.Err=&SyntaxErr{Info:"缺少执行语句end"}
            $$.Err.Scope=$1.Scope
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