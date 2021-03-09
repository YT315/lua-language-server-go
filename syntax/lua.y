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
            if funstmt, ok := $1.(*FuncCall); !ok {
               lualex.(*Lexer).Error("parse error")
            } else {
                $$ = funstmt
                $$.SetLine($1.Line())
            }
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
            $$.End = $6[len(%6)-1].End
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
        TFor TName '=' expr expr TDo block TEnd {
            name:=&NameExpr{Value:$2.Str}
            $$ = &ForLoopNumStmt{Name: name, Init: $4, Limit: $5, Block: $7}
            $$.Start=$1.Start
            $$.End = $7.End
            $$.Err=&SyntaxErr{Info:"缺少逗号"}
            $$.Err.Start=$4.End
            $$.Err.End=$5.Start
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




        /*************** TFunction funcname funcbody *****************/
        TFunction funcname funcbody {
            $$ = $2
            $$.(*FuncDefStmt).Function= $3
            $$.SetLine($3.Line())
            $$.SetLastLine($3.LastLine())
        } |
        /*************** TLocal TFunction TName funcbody *****************/
        TLocal TFunction TName funcbody {
            name:=&NameExpr{Value:$3.Str}
            $$ = &LocalFuncDefStmt{Name: name, Function: $4}
            $$.SetLine($1.line)
            $$.SetLastLine($4.LastLine())
        } | 
        /*************** TLocal namelist '=' exprlist *****************/
        TLocal namelist '=' exprlist {
            $$ = &LocalVarDef{Names: $2, Inits:$4}
            $$.SetLine($1.line)
        } |
        /*************** TLocal namelist *****************/
        TLocal namelist {
            $$ = &LocalVarDef{Names: $2}
            $$.SetLine($1.line)
        }

elseifs: 
        {
            $$ = []Stmt{}
        } | 
        elseifs TElseIf expr TThen block {
            ifstmt:=&IfStmt{Condition: $3, Then: $5}
            $$ = append($1, ifstmt)
            ifstmt.SetLine($2.line)
        }
        
returnstat:
        TReturn {
            $$ = &ReturnStmt{Exprs:nil}
            $$.SetLine($1.line)
        } |
        TReturn exprlist {
            $$ = &ReturnStmt{Exprs:$2}
            $$.SetLine($1.line)
        } |
        TReturn exprlist ';' {
            $$ = &ReturnStmt{Exprs:$2}
            $$.SetLine($1.line)
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
        } |
        funcnameaux ':' TName {
            name:= &NameExpr{Value:$3.Str}
            name.SetLine($3.line)
            $$ = &FuncDefStmt{Name: name, Receiver:$1}
        }

funcnameaux:
        TName {
            $$ = &NameExpr{Value:$1.Str}
            $$.SetLine($1.line)
        } | 
        funcnameaux '.' TName {
            name:= &NameExpr{Value:$3.Str}
            name.SetLine($3.line)
            $$ = &GetItemExpr{Table:$1, Key:name}
            $$.SetLine($3.line)
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
            $$.SetLine($1.line)
        } |
        prefixexp '[' expr ']' {
            $$ = &GetItemExpr{Table:$1, Key:$3}
            $$.SetLine($1.Line())
        } | 
        prefixexp '.' TName {
            name:= &NameExpr{Value:$3.Str}
            name.SetLine($3.line)
            $$ = &GetItemExpr{Table:$1, Key:name}
            $$.SetLine($3.line)
        }

namelist:
        TName {
            name:= &NameExpr{Value:$1.Str}
            name.SetLine($1.line)
            $$ = []Expr{name}
        } | 
        namelist ','  TName {
            name:= &NameExpr{Value:$3.Str}
            name.SetLine($3.line)
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
            $$.SetLine($1.line)
        } | 
        TFalse {
            $$ = &FalseExpr{}
            $$.SetLine($1.line)
        } | 
        TTrue {
            $$ = &TrueExpr{}
            $$.SetLine($1.line)
        } | 
        TNumber {
            $$ = &NumberExpr{Value: $1.Str}
            $$.SetLine($1.line)
        } | 
        TString {
            $$ = &StringExpr{Value: $1.Str}
            $$.SetLine($1.line)
        } |
        TAny {
            $$ = &AnyExpr{}
            $$.SetLine($1.line)
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
            $$.SetLine($1.Line())
        } |
        expr '-' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '*' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '/' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TWdiv expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '^' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '%' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '&' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '~' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '|' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TRmove expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TLmove expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TConn expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '<' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TLequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr '>' expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TBequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TEqual expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TNequal expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TAnd expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        expr TOr expr {
            $$ = &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            $$.SetLine($1.Line())
        } |
        '-' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.SetLine($2.Line())
        } |
        TNot expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.SetLine($2.Line())
        } |
        '#' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.SetLine($2.Line())
        } |
        '~' expr %prec UNARY {
            $$ = &OneOpExpr{Operator: $1.Str, Target: $2}
            $$.SetLine($2.Line())
        } |
        '(' expr ')' {
            $$ = $2
            $$.SetLine($1.line)
        }

prefixexp:
        var {
            $$ = $1
        } |
        functioncall {
            if funcnode, ok := $1.(*FuncCall); !ok {
               lualex.(*Lexer).Error("parse error")
            } else {
              $$ = funcnode
              $$.SetLine($1.Line())
            }
        } 

functioncall:
        prefixexp args {
            $$ = &FuncCall{Function: $1, Args: $2}
            $$.SetLine($1.Line())
        } |
        prefixexp ':' TName args {
            name:= &NameExpr{Value: $3.Str}
            name.SetLine($3.line)
            $$ = &FuncCall{Function: name, Receiver: $1, Args: $4}
            $$.SetLine($1.Line())
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
            str.SetLine($1.line)
            $$ = []Expr{str}
        }

functiondef:
        TFunction funcbody {
            $$ = $2
        }

funcbody:
        '(' parlist ')' block TEnd {
            $$ = &FuncDefExpr{Param: $2, Block: $4}
            $$.SetLine($1.line)
            $$.SetLastLine($5.line)
        } | 
        '(' ')' block TEnd {
            $$ = &FuncDefExpr{Param: nil, Block: $3}
            $$.SetLine($1.line)
            $$.SetLastLine($4.line)
        }

parlist:
        TAny {
            $$ = &ParamExpr{IsAny: true}
            $$.SetLine($1.line)
        } | 
        namelist {
            $$ = &ParamExpr{Params: $1, IsAny: false}
        } | 
        namelist ',' TAny {
            $$ = &ParamExpr{Params: $1, IsAny: true}
        }


tableconstructor:
        '{' '}' {
            $$ = &TableExpr{Fields: []Expr{}}
            $$.SetLine($1.line)
        } |
        '{' fieldlist '}' {
            $$ = &TableExpr{Fields: $2}
            $$.SetLine($1.line)
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
            $$ = &FieldExpr{Key: name, Value: $3}
            $$.SetLine($1.line)
        } | 
        '[' expr ']' '=' expr {
            $$ = &FieldExpr{Key: $2, Value: $5}
            $$.SetLine($2.Line())
        } |
        expr {
            $$ = &FieldExpr{Value: $1}
        }

fieldsep:
        ',' {
            $$ = &NameExpr{Value: $1.Str}
        } | 
        ';' {
            $$ = &NameExpr{Value: $1.Str}
        }
%%