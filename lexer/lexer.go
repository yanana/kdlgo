package lexer

import (
	"fmt"
	"unicode/utf8"
)

type TokenKind int8

const (
	TokUnknown TokenKind = iota - 1
	TokIdentifier
	TokNodeComment
	TokNull TokenKind = iota
	TokString
	TokBool
	TokPlus
	TokComma
	TokSemicolon
	TokEqual
	TokOpenPar
	TokClosePar
	TokOpenBrace
	TokCloseBrace
	TokOpenBracket
	TokCloseBracket
	TokHexadecimal
	TokBinary
	TokOctal
	TokDecimal
	TokFloat
)

type Token struct {
	lexeme string
	start  int
	kind   TokenKind
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

func (l *Lexer) literal(lit string, consume bool) bool {
	if l.until(len(lit)) == lit {
		if consume {
			l.consume(len(lit))
		}

		return true
	}

	return false
}

func (l *Lexer) skipWhile(predicate func(rune) bool) bool {
	before := l.current
	for !l.eof() {
		next := l.peek()
		if !predicate(next) {
			return l.current > before
		}

		l.consume(utf8.RuneLen(next))
	}

	return l.current > before
}

func (l *Lexer) skipChars(skipped []rune) bool {
	return l.skipWhile(func(r rune) bool {
		for _, s := range skipped {
			if r == s {
				return true
			}
		}

		return false
	})
}

var (
	digits         = []rune("0123456789")
	nonIdentifiers = []rune("\\/(){}<>;[]=,\"")
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
		"\u0085", // Next Line
		"\f",
		"\u2028", // Line Separator
		"\u2029", // Paragraph Separator
	}
	literals = map[string]TokenKind{
		"null":  TokNull,
		"true":  TokBool,
		"false": TokBool,
		"+":     TokPlus,
		",":     TokComma,
		";":     TokSemicolon,
		"=":     TokEqual,
		"(":     TokOpenPar,
		")":     TokClosePar,
		"{":     TokOpenBrace,
		"}":     TokCloseBrace,
		"[":     TokOpenBracket,
		"]":     TokCloseBracket,
	}
)

func isNonInitialCharacter(r rune) bool {
	return isNonIdentifier(r) || isDigit(r)
}

func isDigit(r rune) bool {
	for _, d := range digits {
		if d == r {
			return true
		}
	}

	return false
}

func isOctal(r rune) bool {
	return '0' <= r && r <= '7'
}

func isHexadecimal(r rune) bool {
	return isDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func isBinary(r rune) bool {
	return r == '0' || r == '1'
}

func isNonIdentifier(r rune) bool {
	if r <= 0x20 || 0x10FFFF < r {
		return true
	}

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

func (l *Lexer) tokenizeWhiteSpace(consume bool) bool {
	fmt.Printf("source: %s, current: %d, stack: %+v\n", l.source, l.current, l.stack)
	for !l.eof() && isWhiteSpace([]rune(l.source)[l.current]) {
		if consume {
			l.consume(1)
		}

		return true
	}

	return false
}

func (l *Lexer) tokenizeNewLine(consume bool) bool {
	for _, nl := range newLines {
		if l.eofAt(len(nl)) {
			continue
		}
		if l.source[l.current:l.current+len(nl)] == nl {
			if consume {
				l.consume(len(nl))
			}
			return true
		}
	}

	return false
}

func (l *Lexer) tokenizeString(consume bool) /*(l *Lexer)*/ bool {
	before := l.current

	if l.peek() != '"' {
		l.current = before
		return false
	}

	if consume {
		l.consume(1)
	}

	terminated := false

L:
	for !l.eof() {
		switch l.peek() {
		case '\\':
			// TODO: escape sequences
		case '"':
			if consume {
				l.consume(1)
			}
			terminated = true
			break L
		default:
			if consume {
				l.consume(1)
			}
		}
	}

	if !terminated {
		panic("unterminated string")
	}

	l.addToken(TokString, before, l.current)

	return true
}

func (l *Lexer) tokenizeLiteral(consume bool) bool {
	if l.eof() {
		return false
	}

	before := l.current

	for lit, kind := range literals {
		if l.eof() {
			return false
		}
		if l.literal(lit, consume) {
			l.addToken(kind, before, l.current)

			return true
		}
	}

	return false
}

func (l *Lexer) tokenizeNodeComment(consume bool) bool {
	if l.eof() {
		return false
	}

	before := l.current

	if l.literal("/-", consume) {
		l.addToken(TokNodeComment, before, l.current)

		return true
	}

	return false
}

func (l *Lexer) tokenizeIdentifier(consume bool) bool {
	fmt.Printf("char: %q, nonInitial: %t\n", l.peek(), isNonInitialCharacter(l.peek()))
	if l.eof() || isNonInitialCharacter(l.peek()) {
		return false
	}

	before := l.current

	if l.literal("true", true) || l.literal("false", true) || l.literal("null", true) {
		if l.eof() || l.tokenizeWhiteSpace(false) || l.tokenizeNewLine(false) || isNonIdentifier(l.peek()) {
			l.current = before

			return false
		}
	}

L:
	for i := l.current; i < len([]rune(l.source)); i++ {
		r := []rune(l.source)[i]
		if r <= 0x20 || 0x10FFF < r || l.eof() || isWhiteSpace(r) || l.tokenizeNewLine(false) {
			break L
		}
		fmt.Printf("current: %d\n", l.current)
		for _, c := range nonIdentifiers {
			if c == r {
				break L
			}
		}

		fmt.Printf("rune: %c, current: %d\n", r, l.current)
		l.consume(utf8.RuneLen(r))
	}

	if l.current > before {
		l.addToken(TokIdentifier, before, l.current)

		return true
	}

	return false
}

func (l *Lexer) tokenizeHexadecimal(consume bool) bool {
	if l.eof() {
		return false
	}

	before := l.current

	if !l.literal("0x", true) {
		return false
	}

	if l.eof() || !isHexadecimal(l.peek()) {
		l.current = before

		return false
	}

	for !l.eof() {
		next := l.peek()
		if isHexadecimal(next) || next == '_' {
			if consume {
				l.consume(1)
			}
			continue
		}
	}

	if l.current == before {
		return false
	}

	l.addToken(TokHexadecimal, before, l.current)

	return true
}

func (l *Lexer) tokenizeBinary(consume bool) bool {
	if l.eof() {
		return false
	}

	before := l.current

	if !l.literal("0b", true) {
		return false
	}

	// The first character must be either 0 or 1.
	if l.eof() || !isBinary(l.peek()) {
		l.current = before

		return false
	}

	for !l.eof() {
		next := l.peek()
		if isBinary(next) || next == '_' {
			if consume {
				l.consume(1)
			}
			continue
		}
	}

	if l.current == before {
		return false
	}

	l.addToken(TokBinary, before, l.current)

	return true
}

func (l *Lexer) tokenizeOctal(consume bool) bool {
	if l.eof() {
		return false
	}

	before := l.current

	if !l.literal("0o", true) {
		return false
	}

	if l.eof() || !isOctal(l.peek()) {
		l.current = before

		return false
	}

	for !l.eof() {
		next := l.peek()
		if isOctal(next) || next == '_' {
			if consume {
				l.consume(1)
			}
			continue
		}
	}

	if l.current == before {
		return false
	}

	l.addToken(TokOctal, before, l.current)

	return true
}

type lexFn func(l *Lexer, consume bool) bool

var choices = []lexFn{
	(*Lexer).tokenizeWhiteSpace,
	(*Lexer).tokenizeNewLine,
	(*Lexer).tokenizeString,
	(*Lexer).tokenizeLiteral,
	(*Lexer).tokenizeNodeComment,
	(*Lexer).tokenizeIdentifier,
}

func (l *Lexer) Scan() error {
	var matchedSomehow bool
	for !l.eof() {
		matchedSomehow = false

		for _, choice := range choices {
			if choice(l, true) {
				matchedSomehow = true
			}
		}

		if !matchedSomehow {
			return fmt.Errorf("could not match any pattern at %d", l.current)
		}
	}

	return nil
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,
		current: 0,
		stack:   make([]Token, 0),
	}
}
