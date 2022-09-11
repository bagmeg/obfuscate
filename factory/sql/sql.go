package sql

import (
	"bytes"
	"fmt"
	"github.com/bagmeg/obfuscate/config"
)

type SQLObfuscator struct {
	scanner   *SQLScanner
	cfg       *config.SQLConfig
	lastToken TokenKind
}

func NewSQLObfuscator(cfg config.SQLConfig) *SQLObfuscator {
	return &SQLObfuscator{
		scanner: newScanner(&cfg),
		cfg:     &cfg,
	}
}

func (s *SQLObfuscator) Scan(query string) (string, error) {
	s.scanner.Reset(query)

	out := bytes.NewBuffer(make([]byte, 0, len(s.scanner.buf)))
	for {
		token, buf := s.scanner.Scan()

		if token == EndChar {
			break
		}
		if token == LexError {
			return "", fmt.Errorf("%v", s.scanner.Err())
		}
		token, buf = s.Filter(token, buf)
		token, buf = s.consecutiveFilter(token, buf)

		if buf != nil {
			if out.Len() != 0 {
				switch token {
				case ',':
				case SKIPSPACE:
				default:
					if s.lastToken != SKIPNEXTSPACE {
						out.WriteRune(' ')
					}
				}
			}
		}
		out.Write(buf)
		s.lastToken = token
	}
	return out.String(), nil

}

func (s *SQLObfuscator) Filter(token TokenKind, buf []byte) (TokenKind, []byte) {
	switch token {
	case ID:
		if s.lastToken == FROM {
			return TABLE, buf
		}
		return token, buf
	case STRING:
		return s.stringFilter(token, buf)
	case NUMBER:
		return token, questionMark
	default:
		return token, buf
	}
}

func (s *SQLObfuscator) consecutiveFilter(token TokenKind, buf []byte) (TokenKind, []byte) {
	switch token {
	case STRING, NUMBER, POSITIONAL:
		switch s.lastToken {
		case STRING, NUMBER, POSITIONAL:
			return token, nil
		default:
			return token, buf
		}
	default:
		return token, buf
	}
}

func (s *SQLObfuscator) stringFilter(token TokenKind, buf []byte) (TokenKind, []byte) {
	return token, questionMark
}
