package interpreter

import (
	"GLox/lerror"
	"GLox/scanner"
	"fmt"
)

type RuntimeError struct {
	token   *scanner.Token
	message string
}

func NewRuntimeError(token *scanner.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (r *RuntimeError) Error() string {
	lerror.HadRuntimeError = true
	return fmt.Sprintf("Runtime error at line %d : %s", r.token.Line, r.message)
}
