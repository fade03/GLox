package parser

import (
	"LoxGo/lerror"
	"LoxGo/scanner"
	"fmt"
)

type ParseError struct {
	token   *scanner.Token
	message string
}

func newParseError(token *scanner.Token, message string) *ParseError {
	return &ParseError{token: token, message: message}
}

func (e *ParseError) Error() string {
	lerror.HadError = true
	if e.token.Type == scanner.EOF {
		return fmt.Sprintf("[Syntax error] line %d at end %s", e.token.Line, e.message)
	}
	return fmt.Sprintf("[Syntax error] line %d at '%s': %s", e.token.Line, e.token.Lexeme, e.message)
}
