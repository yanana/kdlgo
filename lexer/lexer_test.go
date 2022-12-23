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
				l.tokenizeString()
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
				input: `==`,
				want:  []Token{{kind: TokEqual, lexeme: `==`, start: 0}},
			},
			{
				name:  "not equal",
				input: `!=`,
				want:  []Token{{kind: TokNotEqual, lexeme: `!=`, start: 0}},
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				l := NewLexer(tt.input)
				l.tokenizeLiteral()
				assert.Equal(t, tt.want, l.stack)
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
