package syntax

type Token struct {
	Type rune //对应的tocken做
	Str  string
	Scope
}

type Scope struct {
	Start Pos
	End   Pos
}

type Pos struct {
	line, col uint
}

type SyntaxErr struct {
	Scope
	Info string
}

func (s *SyntaxErr) Error() string {
	return s.Info
}
