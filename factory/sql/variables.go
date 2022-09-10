package sql

import "unicode"

const (
	LexError = TokenKind(65536) + iota
	STRING
	SELECT
	WHERE
	FROM
	ORDER
	BY
	AS
	JOIN
	ID
	NUMBER
	NOT
	TABLE
	COLON
	SKIPSPACE
	SKIPNEXTSPACE
	POSITIONAL
	MARKED
)

const EndChar = unicode.MaxRune + 1

var questionMark = []byte("?")

var tokenKindStrings = map[TokenKind]string{
	STRING: "STRING",
	SELECT: "SELECT",
	WHERE:  "WHERE",
	FROM:   "FROM",
	ORDER:  "ORDER",
	BY:     "BY",
	AS:     "AS",
	JOIN:   "JOIN",
	NUMBER: "NUMBER",
	ID:     "ID",
}

var keywords = map[TokenKind]struct{}{
	SELECT: {},
	WHERE:  {},
	ORDER:  {},
	FROM:   {},
	BY:     {},
	AS:     {},
	JOIN:   {},
	NOT:    {},
}

// TODO keyword 추가 필요
var stringToKeywords = map[string]TokenKind{
	"SELECT": SELECT,
	"WHERE":  WHERE,
	"ORDER":  ORDER,
	"FROM":   FROM,
	"BY":     BY,
	"AS":     AS,
	"JOIN":   JOIN,
	"NOT":    NOT,
}

func (k TokenKind) String() string {
	val, found := tokenKindStrings[k]
	if found {
		return val
	}
	return "<unknown>"
}
