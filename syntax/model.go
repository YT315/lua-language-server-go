package syntax

import "lualsp/protocol"

type Token struct {
	Type rune //对应的tocken做
	Str  string
	Scope
}

type Scope struct {
	Start Pos
	End   Pos
}

//转化为range
func (s Scope) Convert2Range() (res protocol.Range) {
	res.Start.Line = float64(s.Start.line)
	res.Start.Character = float64(s.Start.col)
	res.End.Line = float64(s.End.line)
	res.End.Character = float64(s.End.col)
	return
}

type Pos struct {
	line, col uint
}

type SyntaxErr struct {
	Scope
	Errtype SyntaxErrBase
}

//将此错误插入到lex中
func (s *SyntaxErr) insertInto(llex luaLexer) {
	if lex, ok := llex.(*Lexer); ok {
		diag := protocol.Diagnostic{
			Range:    s.Convert2Range(),
			Severity: protocol.SeverityError,
			Source:   "syntax",
			Message:  string(s.Errtype),
		}
		lex.Diagnostics = append(lex.Diagnostics, diag)
	} else {
		//logger
	}
}

func (s *SyntaxErr) Error() string {
	return string(s.Errtype)
}
