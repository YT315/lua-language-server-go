%{
package syntax

import("strconv")
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

//符号         name   ==      ~=     <=       >=      <<     >>    //     ::     ..   ...  string  number    +    @    
%token<token> TName TEqual TNequal TLequal TBequal TLmove TRmove TWdiv TTarget TConn TAny TString TNumber TAType TSType 
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
%type<expr> name

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
            temp.Start = $1[0].GetStart()
            if len($3)>0{
                temp.End = $3[len($3)-1].GetEnd()
            }else{
                temp.End = $2.End
            }
            
            $$ = temp
        } |
        varlist{
            temp := &AssignStmt{Left: $1, Right: nil}
            temp.Start=$1[0].GetStart()
            temp.End=$1[0].GetEnd()
            temp.Err=&SyntaxErr{Errtype:LackRight}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |  
        varlist '=' {
            temp := &AssignStmt{Left: $1, Right: nil}
            temp.Start=$1[0].GetStart()
            temp.End=$2.End
            temp.Err=&SyntaxErr{Errtype:LackRight}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |  
        '=' exprlist {
            temp := &AssignStmt{Left: nil, Right: $2}
            temp.Start=$1.Start
            if len($2)>0{
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.End
            }
            temp.Err=&SyntaxErr{Errtype:LackLeft}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |
        /*************** functioncall *****************/
        functioncall {
            if funstmt, ok := $1.(*FuncCall); !ok {
               lualex.(*Lexer).Error("parse error")
            } else {
              $$ = funstmt
            }
        } |
        /*************** label *****************/
        label {
            $$ = $1
        } |  
        /*************** TBreak *****************/
        TBreak {
            $$ = &BreakStmt{}
            $$.(*BreakStmt).Scope = $1.Scope
        } |  
        /*************** goto Name *****************/
        TGoto name {
            temp := &GotoStmt{Name:$2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
        } |
        TGoto {
            temp := &GotoStmt{Name:nil}
            temp.Scope = $1.Scope
            temp.Err=&SyntaxErr{Errtype:LackGotoName}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)  
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
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.End
            }

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=temp.Scope
            err.insertInto(lualex)
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
                temp.End = $4[len($4)-1].GetEnd()
            }else{
                temp.End = $3.End
            }

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=$1.Scope
            err.insertInto(lualex)
            $$ = temp
        } |
        TWhile TDo block TEnd {
            temp := &WhileStmt{Condition: nil, Block: $3}
            temp.Start=$1.Start
            temp.End = $4.End
            temp.Err=&SyntaxErr{Errtype:LackWhileCond}
            temp.Err.Start=$1.End
            temp.Err.End=$2.Start     
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |
        TWhile expr {
            temp := &WhileStmt{Condition: $2, Block: nil}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            temp.Err=&SyntaxErr{Errtype:LackWhileBlock}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |
        TWhile {
            temp := &WhileStmt{Condition: nil, Block: nil}
            temp.Scope=$1.Scope
            temp.Err=&SyntaxErr{Errtype:LackWhileCond}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |
        /*************** repeat block until exp *****************/ 
        TRepeat block TUntil expr {
            temp := &RepeatStmt{Condition: $4, Block: $2}
            temp.Start=$1.Start
            temp.End = $4.GetEnd()
            $$ = temp
        } |
        TRepeat block TUntil {
            temp := &RepeatStmt{Condition: nil, Block: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            temp.Err=&SyntaxErr{Errtype:LackUntilCond}
            temp.Err.Scope=$3.Scope
            temp.Err.insertInto(lualex)  
            $$ = temp
        } |
        TRepeat block {
            temp := &RepeatStmt{Condition: nil, Block: $2}
            temp.Start=$1.Start
            if len($2)>0 {
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.End
            }
            temp.Err=&SyntaxErr{Errtype:LackUntilCond}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)  

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
            $$.(*IfStmt).Start=$1.Start
            $$.(*IfStmt).End = $6.End
        } |
        TIf expr TThen block elseifs {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            if len($5)>0{
                temp.End = $5[len($5)-1].GetEnd()
            }else if len($4)>0 {
                temp.End = $4[len($4)-1].GetEnd()
            }else{
                temp.End = $3.End
            }

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=temp.Scope
            err.insertInto(lualex)
        } |
        TIf TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            temp.End = $5.End

            temp.Err=&SyntaxErr{Errtype:LackIfCond}
            temp.Err.Start=$1.End
            temp.Err.End=$2.Start 
            temp.Err.insertInto(lualex)  
        } |
        TIf block elseifs TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            temp.End = $4.End
            
            temp.Err=&SyntaxErr{Errtype:LackIfCond}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

        } |
        //err
        TIf block elseifs {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            if len($3)>0{
                temp.End = $3[len($3)-1].GetEnd()
            }else if len($2)>0{
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End=$1.End
            }

            temp.Err=&SyntaxErr{Errtype:LackIfCond}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)

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
            $$.(*IfStmt).Start = $1.Start
            $$.(*IfStmt).End = $8.End
        } |
        TIf TThen block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $3}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $6
            temp :=$$.(*IfStmt)
            temp.Start = $1.Start
            temp.End = $7.End

            temp.Err=&SyntaxErr{Errtype:LackIfCond}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)

        } |
        TIf expr TThen block elseifs TElse block {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $4 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $7
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            if len($7)>0 {
                temp.End = $7[len($7)-1].GetEnd()
            }else{
                temp.End = $6.End
            }

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=temp.Scope
            err.insertInto(lualex)
           
        } |
        TIf block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: nil, Then: $2}
            cur := $$
            for _, elseif := range $3 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $5
            temp :=$$.(*IfStmt)
            temp.Start=$1.Start
            temp.End = $6.End
            
            temp.Err=&SyntaxErr{Errtype:LackIfCond}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)
        } |
        /*************** TFor TName '=' expr ',' expr TDo block TEnd *****************/
        TFor name '=' expr ',' expr TDo block TEnd {
            $$ = &ForLoopNumStmt{Name: $2, Init: $4, Limit: $6, Block: $8}
            $$.(*ForLoopNumStmt).Start=$1.Start
            $$.(*ForLoopNumStmt).End = $9.End
        } |
        TFor name '=' expr ',' expr TDo block {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: $6, Block: $8}
            temp.Start=$1.Start
            if len($8)>0 {
                temp.End = $8[len($8)-1].GetEnd()
            }else{
                temp.End = $7.End
            }

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=$1.Scope
            err.insertInto(lualex)

            $$ = temp
        } |
        TFor '=' expr ',' expr TDo block TEnd {
            temp := &ForLoopNumStmt{Name: nil, Init: $3, Limit: $5, Block: $7}
            temp.Start=$1.Start
            temp.End = $8.End

            temp.Err=&SyntaxErr{Errtype:LackName}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor TDo block TEnd {
            temp := &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            temp.Start=$1.Start
            temp.End = $4.End

            temp.Err=&SyntaxErr{Errtype:LackForCond}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor TDo block {
            temp := &ForLoopNumStmt{Name: nil, Init: nil, Limit: nil, Block: $3}
            temp.Start=$1.Start
            if len($3)>0 {
                temp.End = $3[len($3)-1].GetEnd()
            }else{
                temp.End = $2.End
            }

            temp.Err=&SyntaxErr{Errtype:LackForCond}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor name TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: nil, Limit: nil, Block: $4}
            temp.Start=$1.Start
            temp.End = $5.End

            temp.Err=&SyntaxErr{Errtype:LackForScope}
            temp.Err.Scope=$2.GetScope()
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor name '=' TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: nil, Limit: nil, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End
            
            temp.Err=&SyntaxErr{Errtype:LackForScope}
            temp.Err.Scope=$2.GetScope()
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor name '=' expr TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: nil, Block: $6}
            temp.Start=$1.Start
            temp.End = $7.End

            temp.Err=&SyntaxErr{Errtype:LackForScope}
            temp.Err.Scope=$4.GetScope()
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor name '=' expr ',' TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: nil, Block: $7}
            temp.Start=$1.Start
            temp.End = $8.End
            temp.Err=&SyntaxErr{Errtype:LackForScope}
            temp.Err.Scope=$5.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        /*************** TFor TName '=' expr ',' expr ',' expr TDo block TEnd *****************/
        TFor name '=' expr ',' expr ',' expr TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: $6, Step:$8, Block: $10}
            temp.Start=$1.Start
            temp.End = $11.End
            $$ = temp
        } |
        TFor name '=' expr ',' expr ',' TDo block TEnd {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: $6, Step:nil, Block: $9}
            temp.Start=$1.Start
            temp.End = $10.End

            temp.Err=&SyntaxErr{Errtype:LackForStep}
            temp.Err.Start=$7.End
            temp.Err.End=$8.Start
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor name '=' expr ',' expr ',' expr TDo block {
            temp := &ForLoopNumStmt{Name: $2, Init: $4, Limit: $6, Step:$8, Block: $10}
            temp.Start=$1.Start
            if len($10)>0{
                temp.End = $10[len($10)-1].GetEnd()
            }else{
                temp.End = $9.End
            }

            temp.Err=&SyntaxErr{Errtype:LackEnd}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        /*************** TFor namelist TIn exprlist TDo block TEnd *****************/
        TFor namelist TIn exprlist TDo block TEnd {
            $$ = &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            $$.(*ForLoopListStmt).Start=$1.Start
            $$.(*ForLoopListStmt).End = $7.End
        } |
        TFor TIn exprlist TDo block TEnd {
            temp := &ForLoopListStmt{Exprs:$3, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End

            temp.Err=&SyntaxErr{Errtype:LackExpr}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor namelist TIn TDo block TEnd {
            temp := &ForLoopListStmt{Names:$2, Block: $5}
            temp.Start=$1.Start
            temp.End = $6.End
            temp.Err=&SyntaxErr{Errtype:LackExpr}
            if len($2)>0{
                temp.Err.Scope=$2[len($2)-1].GetScope()
            }
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } |
        TFor namelist TIn exprlist TDo block {
            temp := &ForLoopListStmt{Names:$2, Exprs:$4, Block: $6}
            temp.Start=$1.Start
            if len($6)>0{
                temp.End = $6[len($6)-1].GetEnd()
            }else{
                temp.End = $5.End
            }
            
            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=$1.Scope
            err.insertInto(lualex)
            
            $$ = temp
        } |
        TFor namelist TIn exprlist{

            temp := &ForLoopListStmt{Names:$2, Exprs:$4}
            temp.Start=$1.Start
            if len($4)>0{
                temp.End = $4[len($4)-1].GetEnd()
            }else{
                temp.End = $3.End
            }
            temp.Err=&SyntaxErr{Errtype:LackBlock}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)

            err:=&SyntaxErr{Errtype:LackEnd}
            err.Scope=$1.Scope
            err.insertInto(lualex)

            $$ = temp
        } |
        /*************** TFunction funcname funcbody *****************/
        TFunction funcname funcbody {
            temp := $2.(*FuncDefStmt)
            temp.Function= $3
            temp.Start=$1.Start
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        TFunction funcname {
            temp := $2.(*FuncDefStmt)
            temp.Start=$1.Start
            temp.Err=&SyntaxErr{Errtype:LackFunctionContent}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        /*TFunction funcbody {
            temp := &FuncDefStmt{}
            temp.Function= $2
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            temp.Err=&SyntaxErr{Errtype:"缺少函数名"}
            temp.Err.Scope=$1.Scope
            $$ = temp
        } |*/
        /*************** TLocal TFunction TName funcbody *****************/
        TLocal TFunction name funcbody {
            temp := &LocalFuncDefStmt{Name: $3, Function: $4}
            temp.Start=$1.Start
            temp.End = $4.GetEnd()
            $$ = temp
        } | 
        TLocal TFunction funcbody {
            temp := &LocalFuncDefStmt{Function: $3}
            temp.Start=$1.Start
            temp.End = $3.GetEnd()

            temp.Err=&SyntaxErr{Errtype:LackFunctionName}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)

            $$ = temp
        } | 
        /* TLocal name funcbody {
            temp := &LocalFuncDefStmt{Name: $2, Function: $3}
            temp.Start=$1.Start
            temp.End = $3.GetEnd()

            err:=&SyntaxErr{Errtype:LackFunctionkeyword}
            err.Scope=$1.Scope
            err.insertInto(lualex)

            $$ = temp
        } |  */
        TLocal TFunction name {
            temp := &LocalFuncDefStmt{Name: $3}
            temp.Start=$1.Start
            temp.End = $3.GetEnd()
            temp.Err=&SyntaxErr{Errtype:LackFunctionContent}
            temp.Err.Scope=$3.GetScope()
            temp.Err.insertInto(lualex)
            $$ = temp
        } | 
        TLocal TFunction {
            temp := &LocalFuncDefStmt{}
            temp.Start=$1.Start
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackFunction}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } | 
        /*************** TLocal namelist '=' exprlist *****************/
        TLocal namelist '=' exprlist {
            temp := &LocalVarDef{Names: $2, Inits:$4}
            temp.Start=$1.Start
            if len($4)>0{
                temp.End = $4[len($4)-1].GetEnd()
            }else{
                temp.End = $3.End
            }
            
            $$ = temp
        } |
        TLocal namelist '=' {
            temp := &LocalVarDef{Names: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            temp.Err=&SyntaxErr{Errtype:LackInitValue}
            temp.Err.Scope=$3.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        TLocal '=' exprlist {
            temp := &LocalVarDef{Inits:$3}
            temp.Start=$1.Start
            if len($3)>0{
                temp.End = $3[len($3)-1].GetEnd()
            }else{
                temp.End = $2.End
            }
            temp.Err=&SyntaxErr{Errtype:LackName}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        /*************** TLocal namelist *****************/
        TLocal namelist {
            temp := &LocalVarDef{Names: $2}
            temp.Start=$1.Start
            if len($2)>0{
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.End
            }
            $$ = temp
        }|
        TLocal {
            temp := &LocalVarDef{}
            temp.Scope=$1.Scope
            temp.Err=&SyntaxErr{Errtype:LackName}
            temp.Err.Scope=$1.Scope
            $$ = temp
        }|
        error{
            temp := &ErrorStmt{Info:"errr"}
            temp.Err=&SyntaxErr{Errtype:StmtErr}
            tk:=lualex.(*Lexer).Token
            temp.Err.Scope=tk.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }

elseifs: 
        {
            $$ = []Stmt{}
        } | 
        elseifs TElseIf expr TThen block {
            ifstmt:=&IfStmt{Condition: $3, Then: $5}
            ifstmt.Start=$2.Start
            if len($5)>0{
                ifstmt.End = $5[len($5)-1].GetEnd()
            }else{
                ifstmt.End = $4.End
            }
            
            $$ = append($1, ifstmt)
        }
        
returnstat:
        TReturn {
            $$ = &ReturnStmt{Exprs:nil}
            $$.(*ReturnStmt).Scope=$1.Scope
        } |
        TReturn exprlist {
            $$ = &ReturnStmt{Exprs:$2}
            $$.(*ReturnStmt).Start=$1.Start
            if len($2)>0{
                $$.(*ReturnStmt).End = $2[len($2)-1].GetEnd()
            }else{
                $$.(*ReturnStmt).End = $1.End
            }
        } |
        TReturn exprlist ';' {
            $$ = &ReturnStmt{Exprs:$2}
            $$.(*ReturnStmt).Start=$1.Start
            $$.(*ReturnStmt).End = $3.End
        }

label:  
        TTarget name TTarget {
            temp := &LabelStmt{Name: $2}
            temp.Start=$1.Start
            temp.End=$3.End
            $$ = temp
        }|
        name TTarget {
            temp := &LabelStmt{Name: $1}
            temp.Start=$1.GetStart()
            temp.End=$2.End

            err:=&SyntaxErr{Errtype:LabelIncomplete}
            err.Scope=$1.GetScope()
            err.insertInto(lualex)

            $$ = temp
        }|
        TTarget name {
            temp := &LabelStmt{Name: $2}
            temp.Start=$1.Start
            temp.End=$2.GetEnd()

            err:=&SyntaxErr{Errtype:LabelIncomplete}
            err.Scope=$2.GetScope()
            err.insertInto(lualex)

            $$ = temp
        }|
        TTarget TTarget {
            temp := &LabelStmt{Name: nil}
            temp.Start=$1.Start
            temp.End=$2.End
            temp.Err=&SyntaxErr{Errtype:LackLabelName}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }|
        TTarget {
            temp := &LabelStmt{Name: nil}
            temp.Scope=$1.Scope
            temp.Err=&SyntaxErr{Errtype:LackLabelName}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)
            //不完整标签
            err:=&SyntaxErr{Errtype:LabelIncomplete}
            err.Scope=temp.Scope
            err.insertInto(lualex)
            $$ = temp
        }
        
funcname: 
        funcnameaux {
            temp := &FuncDefStmt{Name: $1, Receiver: nil}
            temp.Scope=$1.GetScope()
            $$ = temp
        } |
        funcnameaux ':' name {
            temp := &FuncDefStmt{Name: $3, Receiver:$1}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        }|
        funcnameaux ':' {
            temp := &FuncDefStmt{Receiver:$1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackfuncName}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }|
        ':' name {
            temp := &FuncDefStmt{Name: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
            
        }

funcnameaux:
        name {
            $$ = $1
        } | 
        funcnameaux '.' name {
            temp := &GetItemExpr{Table:$1, Key:$3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        funcnameaux '.' {
            temp := &GetItemExpr{Table:$1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        '.' name {
            temp := &GetItemExpr{Key:$2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            temp.Err=&SyntaxErr{Errtype:LackObject}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } 
varlist:
        var {
            $$ = []Expr{$1}
        } | 
        varlist ',' var {
            $$ = append($1, $3)
        } |
        varlist ',' {
            $$ = $1   ////////////////////////////////////////////errrrrrrrrrrrrrrrrr
        } 

var:
        name {
            $$ = $1
        } |
        prefixexp '[' expr ']' {
            temp := &GetItemExpr{Table:$1, Key:$3}
            temp.Start=$1.GetStart()
            temp.End = $4.End
            $$ = temp
        } | 
        prefixexp '[' ']' {
            temp := &GetItemExpr{Table:$1}
            temp.Start=$1.GetStart()
            temp.End = $3.End
            temp.Err=&SyntaxErr{Errtype:LackIndex}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        prefixexp '[' expr {
            temp := &GetItemExpr{Table:$1, Key:$3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            err:=&SyntaxErr{Errtype:LackRightSquareBrackets}
            err.Scope=$3.GetScope()
            err.insertInto(lualex)
            $$ = temp
        } |
        prefixexp '[' {
            temp := &GetItemExpr{Table:$1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackIndex}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        prefixexp '.' TName {
            name := &StringExpr{Value: $3.Str}
            name.Scope=$3.Scope
            temp := &GetItemExpr{Table:$1, Key:name}
            temp.Start=$1.GetStart()
            temp.End = $3.End
            $$ = temp
        } |
        prefixexp '.' {
            temp := &GetItemExpr{Table:$1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }

name:
        TName {
            temp:= &NameExpr{Value:$1.Str}
            temp.Scope=$1.Scope
            $$ = temp
        } | 
        TName TAType{
            class:= &ATypeExpr{Value:$2.Str}
            class.Scope=$2.Scope
            temp:= &NameExpr{Value:$1.Str, Type:class}
            temp.Scope=$1.Scope
            $$ = temp
        } | 
        TName TSType{
            class:= &STypeExpr{Value:$2.Str}
            class.Scope=$2.Scope
            temp:= &NameExpr{Value:$1.Str, Type:class}
            temp.Scope=$1.Scope
            $$ = temp
        }

namelist:
        name {
            $$ = []Expr{$1}
        } | 
        namelist ','  name {
            $$ = append($1, $3)
        } |
        namelist ',' {
            $$ = $1 ///////////////////////////////////////////////////errrrrrrrrrrrr
        }

exprlist:
        expr {
            $$ = []Expr{$1}
        } |
        exprlist ',' expr {
            $$ = append($1, $3)
        } |
        exprlist ',' {
            $$ = $1 ////////////////////////////////////////////////////errrrrrrrrrrrrrrr
        }

expr:
        TNil {
            temp := &NilExpr{}
            temp.Scope =$1.Scope
            $$ = temp
        } | 
        TFalse {
            temp := &FalseExpr{}
            temp.Scope =$1.Scope
            $$ = temp
        } | 
        TTrue {
            temp := &TrueExpr{}
            temp.Scope =$1.Scope
            $$ = temp
        } | 
        TNumber {
            temp := &NumberExpr{}
            temp.Value,_ = strconv.ParseFloat($1.Str, 64)
            temp.Scope =$1.Scope
            $$ = temp
        } | 
        TString {
            temp := &StringExpr{Value: $1.Str}
            temp.Scope =$1.Scope
            $$ = temp
        } |
        TAny {
            temp := &AnyExpr{}
            temp.Scope =$1.Scope
            $$ = temp
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
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '+' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '-' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            
            $$ = temp
        } |
        expr '-' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '*' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '*' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '/' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '/' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TWdiv expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TWdiv {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '^' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '^' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '%' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '%' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '&' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '&' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '~' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '~' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '|' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '|' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TRmove expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TRmove {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TLmove expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TLmove {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TConn expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TConn {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '<' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '<' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TLequal expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TLequal {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr '>' expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr '>' {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TBequal expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TBequal {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TEqual expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TEqual {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TNequal expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TNequal {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TAnd expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TAnd {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        expr TOr expr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1, Right: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } |
        expr TOr {
            temp := &TwoOpExpr{Operator: $2.Str, Left: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackField}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        '-' expr %prec UNARY {
            temp := &OneOpExpr{Operator: $1.Str, Target: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
        } |
        TNot expr %prec UNARY {
            temp := &OneOpExpr{Operator: $1.Str, Target: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
        } |
        '#' expr %prec UNARY {
            temp := &OneOpExpr{Operator: $1.Str, Target: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
        } |
        '~' expr %prec UNARY {
            temp := &OneOpExpr{Operator: $1.Str, Target: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            $$ = temp
        } |
        '(' expr ')' {
            $$ = $2
            $$.setBracket(true)
        } /*|
        '(' expr  {
            temp := $2
            temp.Start=$1.Start
            temp.End = $2.GetEnd()
            temp.Err=&SyntaxErr{Errtype:"缺少右括号"}
            temp.Err.Scope=$1.Scope
        } |
        '(' {
            temp := $2
            temp.Start=$1.Start
            temp.End = $3.End
            temp.Err=&SyntaxErr{Errtype:"缺少右括号"}
            temp.Err.Scope=$2.Scope
        }*/

prefixexp:
        var {
            $$ = $1
        } |
        functioncall {
            if funstmt, ok := $1.(*FuncCall); !ok {
               lualex.(*Lexer).Error("parse error")
            } else {
              $$ = funstmt
            }
        } 

functioncall:
        prefixexp args {
            temp := &FuncCall{Function: $1, Args: $2}
            temp.Start=$1.GetStart()
            if len($2)>0{
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.GetEnd()
            }
            $$ = temp
        } |
        prefixexp ':' name args {
            temp := &FuncCall{Function: $3, Receiver: $1, Args: $4}
            temp.Start=$1.GetStart()
            if len($4)>0{
                temp.End = $4[len($4)-1].GetEnd()
            }else{
                temp.End = $3.GetEnd()
            }
            $$ = temp
        } | 
        prefixexp ':' name  {
            temp := &FuncCall{Function: $3, Receiver: $1,}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            temp.Err=&SyntaxErr{Errtype:LackfuncArgs}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } | 
        prefixexp ':'  {
            temp := &FuncCall{ Receiver: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackfuncName}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }

args:
        '(' ')' {
            $$ = []Expr{}
        } |
        // '(' {
        //     $$ = []Expr{}/////////////////////////////errrrrrrrrrrr
        // } |
        '(' exprlist ')' {
            $$ = $2
        } |
        // '(' exprlist {
        //     $$ = $2 //////////////////////////////////errrrrrrrrrrrrrrrr
        // } |
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
            temp :=  &FuncDefExpr{}
            temp.Scope =$1.Scope
            temp.Err=&SyntaxErr{Errtype:Lackfuncbody}
            temp.Err.Scope=$1.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        }

funcbody:
        '(' parlist ')' block TEnd {
            temp := &FuncDefExpr{Param: $2, Block: $4}
            temp.Start=$1.Start
            temp.End = $5.End
            $$ = temp
        } | 
        '(' parlist ')' TSType block TEnd {
            class:= &STypeExpr{Value:$4.Str}
            class.Scope=$4.Scope
            temp := &FuncDefExpr{Param: $2, Block: $5 ,Result:class}
            temp.Start=$1.Start
            temp.End = $6.End
            $$ = temp
        } | 
        // '(' parlist ')' block {
        //     temp := &FuncDefExpr{Param: $2, Block: $4}
        //     temp.Start=$1.Start
        //     if len($4)>0 {
        //         temp.End = $4[len($4)-1].GetEnd()
        //     }else{
        //         temp.End = $3.End
        //     }
        //     temp.Err=&SyntaxErr{Errtype:"缺少end"}
        //     temp.Err.Scope=temp.Scope
        //     $$ = temp
        // } | 
        // '(' parlist {
        //     temp := &FuncDefExpr{Param: $2}
        //     temp.Start=$1.Start
        //     temp.End = $2.GetEnd()
        //     temp.Err=&SyntaxErr{Errtype:"缺少右括号及函数体"}
        //     temp.Err.Scope=temp.Scope
        //     $$ = temp
        // } | 
        '(' ')' TSType block TEnd {
            class:= &STypeExpr{Value:$3.Str}
            class.Scope=$3.Scope
            temp := &FuncDefExpr{Param: nil, Block: $4 ,Result: class}
            temp.Start=$1.Start
            temp.End = $5.End
            $$ = temp
        } |
        '(' ')' block TEnd {
            temp := &FuncDefExpr{Param: nil, Block: $3}
            temp.Start=$1.Start
            temp.End = $4.End
            $$ = temp
        }

parlist:
        TAny {
            temp := &ParamExpr{IsAny: true}
            temp.Scope =$1.Scope
            $$ = temp
        } | 
        namelist {
            temp := &ParamExpr{Params: $1, IsAny: false}
            temp.Scope =$1[len($1)-1].GetScope()
            $$ = temp
        } | 
        namelist ',' TAny {
            temp := &ParamExpr{Params: $1, IsAny: true}
            if len($1)>0{
                temp.Start=$1[len($1)-1].GetStart()
            }else{
                temp.Start=$2.Start
            }
            
            temp.End = $3.End
            $$ = temp
        }


tableconstructor:
        '{' '}' {
            temp := &TableExpr{Fields: []Expr{}}
            temp.Scope=$1.Scope
            $$ = temp
        } |
        '{' {
            temp := &TableExpr{Fields: []Expr{}}
            temp.Start=$1.Start
            temp.End = $1.End
            err:=&SyntaxErr{Errtype:LackRightCurlyBrackets}
            err.Scope=$1.Scope
            err.insertInto(lualex)
            $$ = temp
        } |
        '{' fieldlist '}' {
            temp := &TableExpr{Fields: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            $$ = temp
        } |
        '{' fieldlist {
            temp := &TableExpr{Fields: $2}
            temp.Start=$1.Start
              if len($2)>0{
                temp.End = $2[len($2)-1].GetEnd()
            }else{
                temp.End = $1.End
            }
            err:=&SyntaxErr{Errtype:LackRightCurlyBrackets}
            err.Scope=$1.Scope
            err.insertInto(lualex)
            $$ = temp
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
        name '=' expr {
            temp := &FieldExpr{Key: $1, Value: $3}
            temp.Start=$1.GetStart()
            temp.End = $3.GetEnd()
            $$ = temp
        } | 
        name '=' {
            temp := &FieldExpr{Key: $1}
            temp.Start=$1.GetStart()
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackRight}
            temp.Err.Scope=$2.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } | 
        '[' expr ']' '=' expr {
            temp := &FieldExpr{Key: $2, Value: $5}
            temp.Start=$1.Start
            temp.End = $5.GetEnd()
            $$ = temp
        } |
        '[' expr ']' '=' {
            temp := &FieldExpr{Key: $2}
            temp.Start=$1.Start
            temp.End = $4.End
            temp.Err=&SyntaxErr{Errtype:LackRight}
            temp.Err.Scope=$3.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        '[' expr ']' {
            temp := &FieldExpr{Key: $2}
            temp.Start=$1.Start
            temp.End = $3.End
            temp.Err=&SyntaxErr{Errtype:LackRight}
            temp.Err.Scope=$3.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        '[' ']' {
            temp := &FieldExpr{}
            temp.Start=$1.Start
            temp.End = $2.End
            temp.Err=&SyntaxErr{Errtype:LackIndex}
            temp.Err.Scope=temp.Scope
            temp.Err.insertInto(lualex)
            $$ = temp
        } |
        '[' {
            temp := &FieldExpr{}
            temp.Start=$1.Start
            temp.End = $1.End
            err:=&SyntaxErr{Errtype:LackRightSquareBrackets}
            err.Scope=temp.Scope
            err.insertInto(lualex)
            $$ = temp
        } |
        '[' expr {
            temp := &FieldExpr{Key: $2}
            temp.Start=$1.Start
            temp.End = $2.GetEnd()

            err:=&SyntaxErr{Errtype:LackRightSquareBrackets}
            err.Scope=$1.Scope
            err.insertInto(lualex)

            $$ = temp
        } |
        expr {
            temp := &FieldExpr{Value: $1}
            temp.Scope =$1.GetScope()
            $$ = temp
        }

fieldsep:
        ',' {
            $$ = &NameExpr{Value: $1.Str}
        } | 
        ';' {
            $$ = &NameExpr{Value: $1.Str}
        }
%%