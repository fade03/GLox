package interpreter

import (
	"GLox/lerror"
	"GLox/scanner/token"
	"fmt"
)

type RuntimeError struct {
	token   *token.Token
	message string
}

func NewRuntimeError(token *token.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (r *RuntimeError) Error() string {
	lerror.HadRuntimeError = true
	return fmt.Sprintf("Runtime error at line %d : %s", r.token.Line, r.message)
}
