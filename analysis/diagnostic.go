package analysis

import (
	"lualsp/protocol"
	"lualsp/syntax"
)

type AnalysisErrBase string

const (
	TypeErr              AnalysisErrBase = "类型错误"
	NotRightValue        AnalysisErrBase = "不能作为右值"
	NoDefine             AnalysisErrBase = "未定义"
	NotTable             AnalysisErrBase = "对象不是表类型"
	IndexErr             AnalysisErrBase = "表索引错误"
	LabelRedef           AnalysisErrBase = "标签重定义"
	LabelFormatErr       AnalysisErrBase = "标签格式错误"
	SyntaxDataErr        AnalysisErrBase = "语法解析错误"
	BreakNoLoopErr       AnalysisErrBase = "break不在循环内"
	LoopValNotNumErr     AnalysisErrBase = "循环值不是数字"
	LoopListParam1Err    AnalysisErrBase = "循环第一个参数必须是函数"
	LoopListParamLoseErr AnalysisErrBase = "循环参数不够"
	LoopListParamMuchErr AnalysisErrBase = "循环参数太多"
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
		Source:   "analysis",
		Message:  string(s.Errtype),
	}
	als.file.Diagnostics = append(als.file.Diagnostics, diag)
}

func (s *AnalysisErr) Error() string {
	return string(s.Errtype)
}
