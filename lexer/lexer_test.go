package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Parallel()

	t.Run("tokenizeString", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
			want  []Token
		}{
			{
				name:  "simple",
				input: `"hello"`,
				want:  []Token{{kind: TokString, lexeme: `"hello"`, start: 0}},
			},
			{
				name:  "sequence",
				input: `"hello" "world"`,
				want:  []Token{{kind: TokString, lexeme: `"hello"`, start: 0}},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				l := NewLexer(tt.input)
				l.tokenizeString(true)
				assert.Equal(t, tt.want, l.stack)
			})
		}
	})

	t.Run("tokenizeLiteral", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
			want  []Token
		}{
			{
				name:  "null",
				input: `null`,
				want:  []Token{{kind: TokNull, lexeme: `null`, start: 0}},
			},
			{
				name:  "false",
				input: `false`,
				want:  []Token{{kind: TokBool, lexeme: `false`, start: 0}},
			},
			{
				name:  "true",
				input: `true`,
				want:  []Token{{kind: TokBool, lexeme: `true`, start: 0}},
			},
			{
				name:  "equal",
				input: `=`,
				want:  []Token{{kind: TokEqual, lexeme: `=`, start: 0}},
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := NewLexer(tt.input)
				l.tokenizeLiteral(true)
				assert.Equal(t, tt.want, l.stack)
			})
		}
	})

	t.Run("tokenizeIdentifier", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
			want  []Token
		}{
			{
				name:  "simple",
				input: `hello`,
				want:  []Token{{kind: TokIdentifier, lexeme: `hello`, start: 0}},
			},
			{
				name:  "literal",
				input: `*`,
				want:  []Token{{kind: TokIdentifier, lexeme: `*`, start: 0}},
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := NewLexer(tt.input)
				l.tokenizeIdentifier(true)
				assert.Equal(t, tt.want, l.stack)
			})
		}
	})

	t.Run("tokenizeBinary", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			input   string
			want    []Token
			matches bool
			current int
		}{
			{
				name:    "simple",
				input:   `+0b1010`,
				want:    []Token{{kind: TokBinary, lexeme: `+0b1010`, start: 0}},
				matches: true,
				current: 7,
			},
			{
				name:    "underscore in the middle",
				input:   `0b101__0_1010`,
				want:    []Token{{kind: TokBinary, lexeme: `0b101__0_1010`, start: 0}},
				matches: true,
				current: 12,
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := NewLexer(tt.input)
				l.tokenizeBinary(true)
				assert.Equal(t, tt.want, l.stack)
			})
		}
	})

	t.Run("tokenizeHexadecimal", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			input   string
			want    []Token
			matches bool
			current int
		}{
			{
				name:    "simple",
				input:   `0x1234`,
				want:    []Token{{kind: TokHexadecimal, lexeme: `0x1234`, start: 0}},
				matches: true,
				current: 6,
			},
			{
				name:    "underscore in the middle",
				input:   `0x1___23_4`,
				want:    []Token{{kind: TokHexadecimal, lexeme: `0x1___23_4`, start: 0}},
				matches: true,
				current: 10,
			},
			{
				name:    "underscore in the middle and at the end",
				input:   `0x1___23_4____`,
				want:    []Token{{kind: TokHexadecimal, lexeme: `0x1___23_4____`, start: 0}},
				matches: true,
				current: 14,
			},
			{
				name:    "underscore at the start",
				input:   `0x_1___23_4`,
				want:    []Token{},
				matches: false,
				current: 0,
			},
			{
				name:    "only prefix",
				input:   `0x`,
				want:    []Token{},
				matches: false,
				current: 0,
			},
			{
				name:    "only prefix and underscore",
				input:   `0x_`,
				want:    []Token{},
				matches: false,
				current: 0,
			},
			{
				name:    "only prefix and underscore",
				input:   `0x_`,
				want:    []Token{},
				matches: false,
				current: 0,
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := NewLexer(tt.input)
				matched := l.tokenizeHexadecimal(true)
				assert.Equal(t, tt.want, l.stack)
				assert.Equal(t, tt.matches, matched)
				assert.Equal(t, tt.current, l.current)
			})
		}
	})

	t.Run("Scan", func(t *testing.T) {
		t.Parallel()

		type fields struct {
			source  string
			current int
			stack   []Token
		}
		tests := []struct {
			name   string
			fields fields
		}{
			{
				name: "simple",
				fields: fields{
					source:  `node`,
					current: 4,
					stack: []Token{
						{kind: TokIdentifier, lexeme: `node`, start: 0},
					},
				},
			},
			{
				name: "string",
				fields: fields{
					source:  `"hello"`,
					current: 7,
					stack: []Token{
						{
							kind:   TokString,
							start:  0,
							lexeme: `"hello"`,
						},
					},
				},
			},
			{
				name: "sequential strings",
				fields: fields{
					source:  `"hello" "world"`,
					current: 15,
					stack: []Token{
						{
							kind:   TokString,
							start:  0,
							lexeme: `"hello"`,
						},
						{
							kind:   TokString,
							start:  8,
							lexeme: `"world"`,
						},
					},
				},
			},
			{
				name: "just a property",
				fields: fields{
					source:  `foo=false`,
					current: 9,
					stack: []Token{
						{
							kind:   TokIdentifier,
							start:  0,
							lexeme: `foo`,
						},
						{
							kind:   TokEqual,
							start:  3,
							lexeme: `=`,
						},
						{
							kind:   TokBool,
							start:  4,
							lexeme: `false`,
						},
					},
				},
			},
			{
				name: "typed string",
				fields: fields{
					source:  `(string)"hello"`,
					current: 15,
					stack: []Token{
						{
							kind:   TokOpenPar,
							start:  0,
							lexeme: "(",
						},
						{
							kind:   TokIdentifier,
							start:  1,
							lexeme: `string`,
						},
						{
							kind:   TokClosePar,
							start:  7,
							lexeme: ")",
						},
						{
							kind:   TokString,
							start:  8,
							lexeme: `"hello"`,
						},
					},
				},
			},
			{
				name: "sequential literals",
				fields: fields{
					source:  `* + ~ ^ , $ ; /- || =`,
					current: 21,
					stack: []Token{
						{kind: TokIdentifier, lexeme: `*`, start: 0},
						{kind: TokPlus, lexeme: `+`, start: 2},
						{kind: TokIdentifier, lexeme: `~`, start: 4},
						{kind: TokIdentifier, lexeme: `^`, start: 6},
						{kind: TokComma, lexeme: `,`, start: 8},
						{kind: TokIdentifier, lexeme: `$`, start: 10},
						{kind: TokSemicolon, lexeme: `;`, start: 12},
						{kind: TokNodeComment, lexeme: `/-`, start: 14},
						{kind: TokIdentifier, lexeme: `||`, start: 17},
						{kind: TokEqual, lexeme: `=`, start: 20},
					},
				},
			},
			{
				name: "sequential mixed",
				fields: fields{
					source:  ` "hello" null  "world"   false`,
					current: 30,
					stack: []Token{
						{
							kind:   TokString,
							start:  1,
							lexeme: `"hello"`,
						},
						{
							kind:   TokNull,
							start:  9,
							lexeme: `null`,
						},
						{
							kind:   TokString,
							start:  15,
							lexeme: `"world"`,
						},
						{
							kind:   TokBool,
							start:  25,
							lexeme: `false`,
						},
					},
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := &Lexer{
					source: tt.fields.source,
				}
				err := l.Scan()
				assert.NoError(t, err)
				assert.Equal(t, tt.fields.current, l.current)
				assert.Equal(t, tt.fields.stack, l.stack)
			})
		}
	})
}
