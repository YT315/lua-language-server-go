package syntax

type Token struct {
	Type      rune //对应的tocken做
	Str       string
	line, col uint
}
