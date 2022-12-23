package lexer

import (
	"fmt"
)

type TokenKind int8

const (
	TokUnknown TokenKind = iota - 1
	TokNull    TokenKind = iota
	TokString
	TokBool
	TokStar
	TokPlus
	TokTilde
	TokCaret
	TokComma
	TokDollar
	TokGreater
	TokSemicolon
	TokSlashDash
	TokDoublePipe
	TokEqual
	TokNotEqual
	TokOpenPar
	TokClosePar
	TokOpenBrace
	TokCloseBrace
	TokOpenBracket
	TokCloseBracket
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

func (l *Lexer) eof() bool {
	return l.eofAt(0)
}

func (l *Lexer) eofAt(offset int) bool {
	return l.current+offset >= len(l.source)
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
	if l.eof() {
		return 0
	}

	return rune(l.source[l.current])
}

func (l *Lexer) until(until int) string {
	if l.eofAt(until - 1) {
		return l.source[l.current:]
	}

	return l.source[l.current : l.current+until]
}

func (l *Lexer) literal(lit string) bool {
	fmt.Println("literal: ", lit)
	if l.until(len(lit)) == lit {
		l.consume(len(lit))

		return true
	}

	return false
}

func (l *Lexer) TokenizeString(raw bool) {
	if raw {
		l.tokenizeRawString()
	} else {
		//l.tokenizeString()
		return
	}
}

var (
	digits         = []rune("0123456789")
	nonIdentifiers = []rune{'\\', '/', '(', ')', '{', '}', '<', '>', ';', '[', ']', '=', ',', '"'}
	whiteSpaces    = []rune{
		'\u0009',
		'\u0020',
		'\u00A0',
		'\u1680',
		'\u2000', '\u2001', '\u2002', '\u2003', '\u2004', '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200A',
		'\u202F',
		'\u205F',
		'\u3000',
	}
	newLines = []string{
		"\r\n",
		"\r",
		"\n",
		"\u0085",
		"\f",
		"\u2028",
		"\u2029",
	}
	literals = map[string]TokenKind{
		"null":  TokNull,
		"true":  TokBool,
		"false": TokBool,
		"*":     TokStar,
		"+":     TokPlus,
		"~":     TokTilde,
		"^":     TokCaret,
		",":     TokComma,
		"$":     TokDollar,
		">":     TokGreater,
		";":     TokSemicolon,
		"/-":    TokSlashDash,
		"||":    TokDoublePipe,
		"==":    TokEqual,
		"!=":    TokNotEqual,
		"(":     TokOpenPar,
		")":     TokClosePar,
		"{":     TokOpenBrace,
		"}":     TokCloseBrace,
		"[":     TokOpenBracket,
		"]":     TokCloseBracket,
	}
)

func isNonIdentifier(r rune) bool {
	for _, ni := range nonIdentifiers {
		if ni == r {
			return true
		}
	}

	return false
}

func isWhiteSpace(r rune) bool {
	for _, ws := range whiteSpaces {
		if ws == r {
			return true
		}
	}

	return false
}

func (l *Lexer) tokenizeWhiteSpace() bool {
	fmt.Printf("source: %s, current: %d\n", l.source, l.current)
	for !l.eof() && isWhiteSpace([]rune(l.source)[l.current]) {
		l.consume(1)
	}

	return false
}

func (l *Lexer) tokenizeNewLine() bool {
	for _, nl := range newLines {
		if l.source[l.current:l.current+len(nl)] == nl {
			l.consume(len(nl))
			return true
		}
	}

	return false
}

func (l *Lexer) tokenizeString() /*(l *Lexer)*/ bool {
	before := l.current

	if l.peek() != '"' {
		l.current = before
		return false
	}

	l.consume(1)

	terminated := false

L:
	for !l.eof() {
		switch l.peek() {
		case '\\':
			// TODO: escape sequences
		case '"':
			l.consume(1)
			terminated = true
			break L
		default:
			l.consume(1)
		}
	}

	if !terminated {
		panic("unterminated string")
	}

	l.addToken(TokString, before, l.current)

	return true
}

func (l *Lexer) tokenizeLiteral() bool {
	if l.eof() {
		return false
	}

	before := l.current

	for lit, kind := range literals {
		if l.eof() {
			break
		}
		if l.literal(lit) {
			l.addToken(kind, before, l.current)
			return true
		}
	}

	return false
}

func (l *Lexer) tokenizeIdentifier() bool {
	if l.eof() || isNonIdentifier(l.peek()) {
		return false
	}

	before := l.current

	for !l.eof() && !isWhiteSpace(l.peek()) {
		l.consume(1)
	}

	l.addToken(TokUnknown, before, l.current)

	return true
}

type lexFn func(l *Lexer) bool

var choices = []lexFn{
	(*Lexer).tokenizeWhiteSpace,
	//(*Lexer).tokenizeNewLine,
	(*Lexer).tokenizeString,
	(*Lexer).tokenizeLiteral,
}

func (l *Lexer) Scan() error {
	var matchedSomehow bool
	for !l.eof() {
		for _, choice := range choices {
			if choice(l) {
				matchedSomehow = true
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
