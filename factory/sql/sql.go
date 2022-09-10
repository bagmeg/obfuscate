package sql

import (
	"github.com/bagmeg/obfuscate/config"
	"log"
)

type SQLObfuscator struct {
	cfg config.SQLConfig
}

func NewSQLObfuscator() *SQLObfuscator {
	return &SQLObfuscator{}
}

func (s *SQLObfuscator) Tokenize() {
	log.Println("SQL Obfuscator Tokenize()")
}

func (s *SQLObfuscator) Parse() {
	log.Println("SQL Obfuscator Parse()")
}
