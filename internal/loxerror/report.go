package loxerror

import (
	"GLox/internal/scanner/token"
	"log"
)

func ReportLexError(line int, where string, message string) {
	HadError = true
	log.Printf("[line %d ] Error %s : %s\n", line, where, message)
}

func ReportResolveError(token *token.Token, message string) {
	HadResolveError = true
	log.Printf("[line %d ] Error %s : %s\n", token.Line, token.Lexeme, message)
}
