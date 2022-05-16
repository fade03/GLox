package scanner

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE
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

	s.addToken(IDENTIFIER, il)
}
