package sql

import (
	"github.com/bagmeg/obfuscate/config"
	"log"
)

type SQLObfuscator struct {
	tokenizer *SQLTokenizer
	cfg       *config.SQLConfig
}

func NewSQLObfuscator(cfg config.SQLConfig) *SQLObfuscator {
	return &SQLObfuscator{
		tokenizer: newTokenizer(&cfg),
		cfg:       &cfg,
	}
}

func (s *SQLObfuscator) Tokenize(query string) {
	s.tokenizer.Reset(query)
	s.tokenizer.Tokenize()
}

func (s *SQLObfuscator) Parse() {
	log.Println("SQL Obfuscator Parse()")
}
