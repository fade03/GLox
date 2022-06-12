package scanner

import (
	"GLox/internal/scanner/token"
)

var keywords map[string]token.TokenType

func init() {
	keywords = make(map[string]token.TokenType)
	keywords["and"] = token.AND
	keywords["class"] = token.CLASS
	keywords["else"] = token.ELSE
	keywords["false"] = token.FALSE
	keywords["for"] = token.FOR
	keywords["fun"] = token.FUN
	keywords["if"] = token.IF
	keywords["nil"] = token.NIL
	keywords["or"] = token.OR
	keywords["print"] = token.PRINT
	keywords["return"] = token.RETURN
	keywords["super"] = token.SUPER
	keywords["this"] = token.THIS
	keywords["true"] = token.TRUE
	keywords["var"] = token.VAR
	keywords["while"] = token.WHILE
}

func (s *Scanner) addIdentifier() {
	for s.isAlphaDigit(s.peek()) {
		s.advance()
	}
	// 判断扫描出的identifier是否是keyword
	il := s.source[s.start:s.current]
	if kw, exits := keywords[il]; exits {
		s.addToken(kw, nil)
		return
	}

	s.addToken(token.IDENTIFIER, il)
}
