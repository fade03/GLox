package token

import "fmt"

const (
	LEFT_PAREN  = iota // '('
	RIGHT_PAREN        // ')'
	LEFT_BRACE         // '{'
	RIGHT_BRACE        // '}'
	COMMA              // ','
	DOT                // '.'
	MINUS              // '-'
	PLUS               // '+'
	SEMICOLON          // ';'
	SLASH              // '/'
	STAR               // '*'

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	IDENTIFIER
	STRING
	NUMBER

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type TokenType = int

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%d %s %v", t.Type, t.Lexeme, t.Literal)
}
