package resolver

import (
	"GLox/internal/scanner/token"
)

type Scope map[string]bool

func (r *Resolver) beginScope() {
	r.scopes.Push(make(Scope))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(token *token.Token) {
	if r.scopes.isEmpty() {
		return
	}

	r.scopes.Peek().(Scope)[token.Lexeme] = false
}

func (r *Resolver) define(token *token.Token) {
	if r.scopes.isEmpty() {
		return
	}

	r.scopes.Peek().(Scope)[token.Lexeme] = true
}
