package sql

import (
	"bytes"
	"fmt"
	"github.com/bagmeg/obfuscate/config"
	"unicode"
	"unicode/utf8"
)

type TokenKind uint32

type token struct {
	token TokenKind
}

type SQLTokenizer struct {
	pos         int
	off         int
	lastChar    rune
	buf         []byte
	err         error
	latestToken TokenKind
	cfg         *config.SQLConfig
	tokens      []token
}

func (tkn *SQLTokenizer) Tokenize() {

}

func newTokenizer(cfg *config.SQLConfig) *SQLTokenizer {
	if cfg == nil {
		cfg = new(config.SQLConfig)
	}
	return &SQLTokenizer{
		cfg: cfg,
	}
}

func (tkn *SQLTokenizer) Reset(in string) {
	tkn.reset(in)
}

func (tkn *SQLTokenizer) reset(in string) {
	tkn.pos = 0
	tkn.off = 0
	tkn.lastChar = 0
	tkn.buf = []byte(in)
	tkn.err = nil
	tkn.tokens = nil
}

func (tkn *SQLTokenizer) Err() error {
	return tkn.err
}

func (tkn *SQLTokenizer) setErr(format string, args ...interface{}) {
	if tkn.err != nil {
		return
	}
	tkn.err = fmt.Errorf("sqltokenizer at position %d: %v", tkn.pos, fmt.Errorf(format, args...))
}

// Filter filters the token according to the config file
func (tkn *SQLTokenizer) Filter() {

}

// Scan scans the given query string and generates token
func (tkn *SQLTokenizer) Scan() (TokenKind, []byte) {
	if tkn.lastChar == 0 {
		tkn.advance()
	}
	tkn.skipBlank()

	switch ch := tkn.lastChar; {
	case isLetter(ch):
		return tkn.scanIdentifier()
	case tkn.cfg.ReplaceDigits && isDigit(ch):
		return tkn.scanNumber()
	default:
		tkn.advance()
		switch ch {
		case EndChar:
			return EndChar, nil
		case '-':
			if isDigit(tkn.lastChar) {
				return tkn.scanNumber()
			}
		case '\'':
			return tkn.scanString('\'')
		case '"':
			fallthrough
		case '`':
			tkn.bytes()
			return SKIPSPACE, []byte("")
		case '<':
		case '>':
			return TokenKind(ch), tkn.bytes()
		case ':':
			tkn.advance()
			return COLON, []byte("::")
		case '$':
			// 숫자인 경우 읽어서 뭉땡이로 취급
			if isDigit(tkn.lastChar) {
				tkn.scanNumber()
				return POSITIONAL, []byte("?")
			}
			return SKIPNEXTSPACE, tkn.bytes()
		case '(':
			// mysql (?) 인 경우 ? 로 처리 필요
			if tkn.lastChar == '?' {
				tkn.advance()
				if tkn.lastChar == ')' {
					tkn.advance()
					tkn.bytes()
					return MARKED, questionMark
				}
			}
			fallthrough
		default:
			return TokenKind(ch), tkn.bytes()
		}
	}
	return LexError, tkn.bytes()
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (tkn *SQLTokenizer) scanIdentifier() (TokenKind, []byte) {
	tkn.advance()
	for isLetter(tkn.lastChar) || isDigit(tkn.lastChar) || tkn.lastChar == '.' || tkn.lastChar == '*' || tkn.lastChar == '_' {
		tkn.advance()
	}
	t := tkn.bytes()
	upper := bytes.ToUpper(t)
	keywordVal, found := stringToKeywords[string(upper[:])]
	if found {
		return keywordVal, t
	}
	return ID, t
}

func (tkn *SQLTokenizer) skipBlank() {
	for {
		if unicode.IsSpace(tkn.lastChar) {
			tkn.advance()
		} else {
			break
		}
	}
	tkn.bytes()
}

func (tkn *SQLTokenizer) scanString(delim rune) (TokenKind, []byte) {
	for {
		// TODO 현재는 delim 안에서  escape character를 무시하고 있지만 처리가 필요한지 확인 필요
		ch := tkn.lastChar
		tkn.advance()
		if ch == delim || ch == '\\' {
			if tkn.lastChar == delim {
				// delimiter 두 개 오는 경우 예) 'he''lo' -> he'lo 가 되야한다.
				// select age as 'AG''E' from table -> select age as "AG'E" from table과 같다
				tkn.advance()
			} else {
				break
			}
		}
		// delim 나오기 전에 query가 끝나는 경우
		if ch == EndChar {
			tkn.setErr("EOF while scanning string")
			return LexError, tkn.bytes()
		}
	}
	return STRING, tkn.bytes()
}

// scanNumber scans consecutive numbers ex) 1234 -> 1234
func (tkn *SQLTokenizer) scanNumber() (TokenKind, []byte) {
	// TODO 8, 16진수도 scan할 수 있게 수정 필요
	// 일단 10진수만  + 소수점
	tkn.scanByBase(10)
	if tkn.lastChar == '.' {
		// base 에 따라 읽는 함수 필요
		tkn.advance()
	}
	tkn.scanByBase(10)
	return NUMBER, tkn.bytes()
}

func (tkn *SQLTokenizer) scanByBase(base int) {
	// if value < base then advance
	// else return

	// base = 10
	// space, ',' 경우 10 보다 작기 때문에 음수값이 나온다....
	// TODO base 보다 작은 ascii code에 대한 작업 처리
	for '0' <= tkn.lastChar && tkn.lastChar <= '9' && int(tkn.lastChar)-'0' < base {
		tkn.advance()
	}
	// TODO base=8, base=16인 경우도 처리 가능하게 수정 필요
}

// next reads next rune
func (tkn *SQLTokenizer) advance() {
	ch, n := utf8.DecodeRune(tkn.buf[tkn.off:])
	if ch == utf8.RuneError {
		tkn.lastChar = EndChar
		if n == 1 {
			tkn.setErr("Invalid encoding: %v", tkn.buf[tkn.off])

		}
		tkn.pos++
		return
	}
	tkn.off += n
	tkn.lastChar = ch
}

func (tkn *SQLTokenizer) bytes() []byte {
	if tkn.lastChar == EndChar {
		ret := tkn.buf[:tkn.off]
		tkn.buf = tkn.buf[tkn.off:]
		tkn.off = 0
		return ret
	}
	lastLen := utf8.RuneLen(tkn.lastChar)
	if lastLen == -1 {
		tkn.setErr("Invalid rune: %v", tkn.lastChar)
		return []byte("")
	}
	ret := tkn.buf[:tkn.off-lastLen]
	tkn.buf = tkn.buf[tkn.off-lastLen:]
	tkn.off = lastLen
	return ret
}
