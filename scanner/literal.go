package scanner

import "strconv"

// addStrLiteral 获取source中的字符串字面量
func (s *Scanner) addStrLiteral() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		serror(s.line, "Unterminated string.")
		return
	}
	// 注意这里要consume掉最后一个 " 号
	s.advance()
	// 实际的字符串字面量要去掉左右的 " 号，所以是start+1:current-1
	s.addToken(STRING, s.source[s.start+1:s.current-1])
}

func (s *Scanner) addNumberLiteral() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	// 小数点不能出现在数字字面量的最后一位
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume掉 '.'
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	literal, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(NUMBER, literal)
}
