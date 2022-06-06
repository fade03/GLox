package loxerror

import (
	"GLox/scanner/token"
	"fmt"
)

var (
	HadError        = false
	HadRuntimeError = false
)

type ParseError struct {
	token   *token.Token
	message string
}

func NewParseError(token *token.Token, message string) *ParseError {
	return &ParseError{token: token, message: message}
}

func (e *ParseError) Error() string {
	HadError = true
	if e.token.Type == token.EOF {
		return fmt.Sprintf("line %d at end %s", e.token.Line, e.message)
	}
	return fmt.Sprintf("line %d at '%s': %s", e.token.Line, e.token.Lexeme, e.message)
}

// #########################

type RuntimeError struct {
	token   *token.Token
	message string
}

func NewRuntimeError(token *token.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (r *RuntimeError) Error() string {
	HadRuntimeError = true
	return fmt.Sprintf("Runtime error at line %d : %s", r.token.Line, r.message)
}
