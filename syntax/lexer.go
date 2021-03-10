package syntax

import (
	"fmt"
	"io"
)

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isLetter(ch rune) bool  { return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

//EOF 终止
const EOF = -1

var keywordMap = map[string]rune{
	"and": TAnd, "break": TBreak, "do": TDo, "else": TElse, "elseif": TElseIf,
	"end": TEnd, "false": TFalse, "for": TFor, "function": TFunction, "goto": TGoto,
	"if": TIf, "in": TIn, "local": TLocal, "nil": TNil, "not": TNot, "or": TOr,
	"return": TReturn, "repeat": TRepeat, "then": TThen, "true": TTrue,
	"until": TUntil, "while": TWhile}

type Lexer struct {
	source
	Block         []Stmt
	Token         Token
	PrevTokenType int
}

func NewLexer(in io.Reader, errh func(line, col uint, msg string)) *Lexer {
	lx := &Lexer{}
	lx.source.init(in, errh)
	return lx
}
func (lx *Lexer) Parse() int {
	return luaParse(lx)
}
func (lx *Lexer) Lex(lval *luaSymType) int {
redo:
	tok := &lx.Token
	// skip white space
	lx.stop()
	for lx.ch == ' ' || lx.ch == '\t' || lx.ch == '\n' || lx.ch == '\r' {
		lx.nextch()
	}
	// token start
	tok.Start.line, tok.Start.col = lx.pos()
	tok.End = tok.Start
	lx.start()

	if isLetter(lx.ch) {
		lx.nextch()
		lx.scanIdent()
		goto finally
	}
	if isDecimal(lx.ch) {
		lx.scanNum()
		goto finally
	}
	switch lx.ch {
	case EOF:
		tok.Type = EOF
	case '-':

		if lx.nextch(); lx.ch == '-' {
			lx.lineComment()
			goto redo
		} else {
			tok.Type = '-'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '"', '\'':
		tok.Type = TString
		lx.scanString()
	case '[':
		if lx.nextch(); lx.ch == '[' || lx.ch == '=' {
			lx.scanMultilineString()
		} else {
			tok.Type = '['
			tok.Str = string(lx.ch)
			tok.End.col++
		}
	case '=':
		if lx.nextch(); lx.ch == '=' {
			tok.Type = TEqual
			tok.Str = "=="
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = '='
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '~':
		if lx.nextch(); lx.ch == '=' {
			tok.Type = TNequal
			tok.Str = "~="
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = '~'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '<':
		if lx.nextch(); lx.ch == '=' {
			tok.Type = TLequal
			tok.Str = "<="
			tok.End.col += 2
			lx.nextch()
		} else if lx.ch == '<' {
			tok.Type = TLmove
			tok.Str = "<<"
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = '<'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '>':
		if lx.nextch(); lx.ch == '=' {
			tok.Type = TBequal
			tok.Str = ">="
			tok.End.col += 2
			lx.nextch()
		} else if lx.ch == '>' {
			tok.Type = TRmove
			tok.Str = ">="
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = '>'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '/':
		if lx.nextch(); lx.ch == '/' {
			tok.Type = TWdiv
			tok.Str = "//"
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = '/'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case ':':
		if lx.nextch(); lx.ch == ':' {
			tok.Type = TTarget
			tok.Str = "::"
			tok.End.col += 2
			lx.nextch()
		} else {
			tok.Type = ':'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '.':
		if lx.nextch(); isDecimal(lx.ch) {
			lx.scanNum()
		} else if lx.ch == '.' {
			if lx.nextch(); lx.ch == '.' {
				tok.Type = TAny
				tok.Str = "..."
				tok.End.col += 3
				lx.nextch()
			} else {
				tok.Type = TConn
				tok.Str = ".."
				tok.End.col += 2
			}
		} else {
			tok.Type = '.'
			tok.Str = string(tok.Type)
			tok.End.col++
		}
	case '+', '*', '%', '^', '#', '&', '|', '(', ')', '{', '}', ']', ';', ',':
		tok.Type = lx.ch
		tok.Str = string(tok.Type)
		tok.End.col++
		lx.nextch()
	default:
		lx.Error("Invalid token")
		goto finally
	}

finally:
	lval.token = lx.Token
	return int(tok.Type)
}

func (lx *Lexer) Error(message string) {
	fmt.Println(message)
	//panic(message)
}

func (lx *Lexer) scanNum() {
	if lx.ch == '0' { // octal
		lx.nextch()
		if lx.ch == 'x' || lx.ch == 'X' {
			lx.nextch()
			for isHex(lx.ch) {
				lx.nextch()
			}
		} else if lx.ch == '.' {
			lx.nextch()
			for isDecimal(lx.ch) {
				lx.nextch()
			}
		}
	} else {
		for isDecimal(lx.ch) {
			lx.nextch()
		}
		if lx.ch == '.' {
			for isDecimal(lx.ch) {
				lx.nextch()
			}
		}
	}
	lit := lx.segment()
	lx.Token.Str = string(lit)
	lx.Token.Type = TNumber
	lx.Token.End.col += uint(len(lit))
}

func (lx *Lexer) scanIdent() {
	// accelerate common case (7bit ASCII)
	for isLetter(lx.ch) || isDecimal(lx.ch) {
		lx.nextch()
	}

	// possibly a keyword
	lit := lx.segment()
	if len(lit) >= 2 {
		if tok, ok := keywordMap[string(lit)]; ok {
			lx.Token.Str = string(lit)
			lx.Token.Type = tok
			lx.Token.End.col += uint(len(lit))
			return
		}
	}
	lx.Token.Str = string(lit)
	lx.Token.Type = TName
	lx.Token.End.col += uint(len(lit))
	return
}

func (lx *Lexer) lineComment() {
	// directive text
	lx.skipLine()
	//s.comment(string(s.segment()))
}

func (lx *Lexer) skipLine() {
	// don't consume '\n' - needed for nlsemi logic
	for lx.ch >= 0 && lx.ch != '\n' {
		lx.nextch()
	}
}

func (lx *Lexer) scanString() {
	start := lx.ch
	lx.nextch()
	for start != lx.ch {
		if lx.ch == '\n' || lx.ch == '\r' || lx.ch < 0 {
			lx.Error("unterminated string")
			return
		}
		/*if lx.ch == '\\' {
			if err := lx.scanEscape(ch, buf); err != nil {
				return err
			}
		} else {
			writeChar(buf, ch)
		}*/
		lx.nextch()
	}
	lit := lx.segment()
	lx.Token.Str = string(lit[1:])
	lx.nextch()
	lx.Token.Type = TString
	lx.Token.End.col += uint(len(lit))
}
func (lx *Lexer) scanMultilineString() {
	start := lx.ch
	lx.nextch()
	for start != lx.ch {
		/*if lx.ch == '\\' {
			if err := lx.scanEscape(ch, buf); err != nil {
				return err
			}
		} else {
			writeChar(buf, ch)
		}*/
		lx.nextch()
	}
	lx.Token.Str = string(lx.segment())
	lx.Token.Type = TString

}
