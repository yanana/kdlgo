package lexer

import "fmt"

type TokenKind int8

const (
	TokUnknown TokenKind = iota - 1
	TokNull    TokenKind = iota
	TokString
)

type Token struct {
	lexeme string
	start  int
	kind   TokenKind
}

func newToken(kind TokenKind, start, until int) Token {
	return Token{
		kind:   kind,
		start:  start,
		lexeme: "",
	}
}

type Lexer struct {
	source  string
	current int
	stack   []Token
}

func (l *Lexer) EOF() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) addToken(kind TokenKind, start, until int) {
	l.stack = append(l.stack, Token{
		kind:   kind,
		start:  start,
		lexeme: l.source[start:until],
	})
}

func (l *Lexer) consume(amount int) {
	l.current += amount
}

func (l *Lexer) peek() rune {
	if l.EOF() {
		return 0
	}

	return rune(l.source[l.current])
}

func (l *Lexer) until(until int) string {
	return l.source[l.current : l.current+until]
}

func (l *Lexer) TokenizeString(raw bool) {
	if raw {
		l.tokenizeRawString()
	} else {
		//l.tokenizeString()
		return
	}
}

var whiteSpaces = []rune{
	'\u0009',
	'\u0020',
	'\u00A0',
	'\u1680',
	'\u2000', '\u2001', '\u2002', '\u2003', '\u2004', '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200A',
	'\u202F',
	'\u205F',
	'\u3000',
}

var newLines = []string{
	"\r\n",
	"\r",
	"\n",
	"\u0085",
	"\f",
	"\u2028",
	"\u2029",
}

func isWhiteSpace(r rune) bool {
	for _, ws := range whiteSpaces {
		if ws == r {
			return true
		}
	}

	return false
}

func lexWhiteSpace(l *Lexer) bool {
	for !l.EOF() && isWhiteSpace([]rune(l.source)[l.current]) {
		l.consume(1)
	}

	return false
}

func lexNewLine(l *Lexer) bool {
	for _, nl := range newLines {
		if l.source[l.current:l.current+len(nl)] == nl {
			l.consume(len(nl))
			return true
		}
	}

	return false
}

func (l *Lexer) lexString() /*(l *Lexer)*/ bool {
	before := l.current
	if l.peek() != '"' {
		l.current = before
		return false
	}

	l.consume(1)

	terminated := false

	for !l.EOF() {
		switch l.peek() {
		case '\\':
			// TODO: escape sequences
		case '"':
			l.consume(1)
			terminated = true
			break
		default:
			l.consume(1)
		}
	}

	if !terminated {
		panic("unterminated string")
	}

	l.addToken(TokString, before, l.current)

	return false
}

type lexFn func(l *Lexer) bool

var choices = []lexFn{
	//lexWhiteSpace,
	//lexNewLine,
	(*Lexer).lexString,
}

func (l *Lexer) Scan() error {
	var matchedSomehow bool
	for !l.EOF() {
		for _, choice := range choices {
			if choice(l) {
				matchedSomehow = true
				break
			}
		}
		if !matchedSomehow {
			return fmt.Errorf("could not match any pattern at %d", l.current)
		}
	}

	return nil
}

func (l *Lexer) tokenizeRawString() {

}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,
		current: 0,
		stack:   make([]Token, 0),
	}
}
