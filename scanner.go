package json_sql_query

import (
	"bytes"
	"strings"
)

type Scanner struct {
	pos   int
	query string
}

var eof = rune(0)

func NewScanner(r string) *Scanner {
	return &Scanner{query: r}
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch {
	case ch == eof:
		return EOF, ""
	case ch == '\'' || ch == '"':
		s.unread()
		return s.readString()
	case isDigit(ch):
		s.unread()
		return s.readNumber()
	case ch == '*':
		return ASTERISK, string(ch)
	case ch == ',':
		return COMMA, string(ch)
	case ch == '<':
		ch := s.read()
		if ch == '=' {
			return LTE, "<="
		} else if ch == '>' {
			return NEQ, "<>"
		} else {
			s.unread()
			return LT, "<"
		}
	case ch == '>':
		ch := s.read()
		if ch == '=' {
			return GTE, ">="
		} else {
			s.unread()
			return GT, ">"
		}
	case ch == '=':
		ch := s.read()
		if ch == '~' {
			return MATCH, "=~"
		} else {
			s.unread()
			return EQ, "="
		}
	case ch == '(':
		return LEFTC, "("
	case ch == ')':
		return RIGTHC, ")"

	}
	return ILLEGAL, string(ch)
}

func (s *Scanner) read() rune {

	if s.pos < len(s.query) {
		ch := rune(s.query[s.pos])
		s.pos++
		return ch
	}
	return eof
}

func (s *Scanner) unread() {
	s.pos--
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '.' && ch != '/' {
			if ch == '-' {
				ch = s.read()
				if ch == '>' {
					_, _ = buf.WriteRune('-')
					_, _ = buf.WriteRune('>')
				} else {
					s.unread()
					s.unread()
					break
				}
			} else if ch == '`' {
				for {
					ch = s.read()
					if ch == '`' {
						break
					}
					_, _ = buf.WriteRune(ch)
				}

			} else {
				if strings.ToUpper(buf.String()) == "NOT" && ch == ' ' {
					_, _ = buf.WriteRune(ch)
				} else {
					s.unread()
					break
				}
			}

		} else {
			if ch != '`' {
				_, _ = buf.WriteRune(ch)
			}
		}
	}

	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	case "WHERE":
		return WHERE, buf.String()
	case "AS":
		return AS, buf.String()
	case "IN":
		return IN, buf.String()
	case "NOT IN":
		return NOTIN, buf.String()
	case "LIKE":
		return LIKE, buf.String()
	case "AND":
		return AND, buf.String()
	case "OR":
		return OR, buf.String()
	}
	return IDENT, buf.String()
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WS, buf.String()
}

func (s *Scanner) readNumber() (tok Token, lit string) {
	typ := INTEGER
	var buf bytes.Buffer

	ch := s.read()
	if ch == eof {
		return EOF, ""
	}
	s.unread()
	for {
		ch = s.read()
		if ch == eof || isWhitespace(ch) {
			break
		}
		if ch != eof && !isDigit(ch) && ch != '.' {
			break
		}
		if ch == '.' {
			typ = FLOAT
		}
		buf.WriteRune(ch)
	}
	s.unread()
	return typ, string(buf.Bytes())
}

func (s *Scanner) readString() (tok Token, lit string) {

	ch := s.read()
	if ch == eof || (ch != '\'' && ch != '"') {
		return EOF, ""
	}

	var buf bytes.Buffer

	for {
		chStr := s.read()
		if chStr == ch {
			return STRING, string(buf.Bytes())
		} else {
			if chStr == eof {
				return UNKNOWN_TYPE, string(buf.Bytes())
			} else {
				buf.WriteRune(chStr)
			}
		}
	}
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }
