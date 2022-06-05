package scanner

import (
	"GLox/loxerror"
	"GLox/scanner/token"
	"GLox/utils"
)

type Scanner struct {
	source  string
	tokens  []*token.Token
	start   int // start指向被扫描词素的第一个字符
	current int // current指向当前处理的字符
	line    int // line指向当前行数
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, line: 1}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		// 下一轮扫描的开始位置就是上一轮扫描的结束位置
		s.start = s.current
		s.scanToken()
	}

	// After scanning source, add EOF to tokens
	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	switch c := s.advance(); c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	// Look ahead 一个字符
	case '!':
		s.addToken(utils.Ternary(s.matchNext('='), token.BANG_EQUAL, token.BANG), nil)
		//s.addToken(TokenType((utils.Ternary(s.matchNext('='), BANG_EQUAL, BANG)).(int)), nil)
	case '=':
		s.addToken(utils.Ternary(s.matchNext('='), token.EQUAL_EQUAL, token.EQUAL), nil)
	case '<':
		s.addToken(utils.Ternary(s.matchNext('='), token.LESS_EQUAL, token.LESS), nil)
	case '>':
		s.addToken(utils.Ternary(s.matchNext('='), token.GREATER_EQUAL, token.GREATER), nil)
	// '/' 需要特殊处理，因为注释也是以 '/' 开头
	case '/':
		if s.matchNext('/') {
			// 获取当前current指向的字符，如果不是换行或者到达文件末尾，则直接consume
			for s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	// 字符串以 '"' 开头
	case '"':
		s.addStrLiteral()
	default:
		if s.isDigit(c) {
			s.addNumberLiteral()
		} else if s.isAlpha(c) {
			// 以字母或下划线开头的被认为是一个identifier
			// 假设匹配到的全是identifier，之后再和keyword区分（最长匹配原则）
			s.addIdentifier()
		} else {
			loxerror.Report(s.line, "", "Unexpected character "+string(c))
		}
	}
}

func (s *Scanner) addToken(tokenType token.TokenType, literal interface{}) {
	s.tokens = append(s.tokens, token.NewToken(tokenType, s.source[s.start:s.current], literal, s.line))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance 获取当前current指向的字符，并将current加一
func (s *Scanner) advance() byte {
	defer func() { s.current++ }()
	return s.source[s.current]
}

// previous 获取当前current-1指向的字符
func (s *Scanner) previous() byte {
	return s.source[s.current-1]
}

// peek 仅获取current指向的字符
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\n'
	}
	return s.source[s.current]
}

// peek 获取current+1指向的字符
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\n'
	}
	return s.source[s.current+1]
}

func (s *Scanner) matchNext(excepted byte) bool {
	if s.isAtEnd() {
		return false
	}
	// 如果相等，则consume掉一个字符
	if s.source[s.current] == excepted {
		s.current++
		return true
	}
	// 如果不相等，则不consume字符，接着进入下一轮扫描
	return false
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (s *Scanner) isAlphaDigit(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
