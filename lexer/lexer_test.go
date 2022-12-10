package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("lexString", func(t *testing.T) {
		input := `"hello"`
		l := NewLexer(input)
		l.lexString()
		assert.Equal(t, l.stack, []Token{{kind: TokString, lexeme: `"hello"`, start: 0}})
	})

	t.Run("Scan", func(t *testing.T) {
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := &Lexer{
					source: tt.fields.source,
				}
				l.Scan()
				assert.Equal(t, tt.fields.current, l.current)
				assert.Equal(t, tt.fields.stack, l.stack)
			})
		}
	})
}
