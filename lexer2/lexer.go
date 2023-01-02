package lexer2

type TokenKind int

const (
	TokError TokenKind = iota
	TokIdentifier
	TokNodeComment
	TokNull
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

const eof = -1

var (
	digits = "0123456789"
)

type Token struct {
	kind   TokenKind // The kind of the token.
	start  int       // The starting position of the token in the input string.
	line   int       // The line number of the token in the input string.
	lexeme string    // The actual string of the token.
}

type Lexer struct {
	input  string
	pos    int
	tokens chan Token
}
