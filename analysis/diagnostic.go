package analysis

import (
	"lualsp/protocol"
	"lualsp/syntax"
)

type AnalysisErrBase string

const (
	TypeErr AnalysisErrBase = "类型错误" //ErrorStmt
)

type AnalysisErr struct {
	syntax.Scope
	Errtype AnalysisErrBase
}

//将此错误插入到lex中
func (s *AnalysisErr) insertInto(als *Analysis) {
	diag := protocol.Diagnostic{
		Range:    s.Convert2Range(),
		Severity: protocol.SeverityError,
		Source:   "syntax",
		Message:  string(s.Errtype),
	}
	als.file.Diagnostics = append(als.file.Diagnostics, diag)
}

func (s *AnalysisErr) Error() string {
	return string(s.Errtype)
}
