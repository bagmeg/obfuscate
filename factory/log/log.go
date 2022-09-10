package log

import (
	"github.com/bagmeg/obfuscate/config"
	"log"
)

type LogObfuscator struct {
	cfg config.LogConfig
}

func NewLogObfuscator() *LogObfuscator {
	return &LogObfuscator{}
}

func (l *LogObfuscator) Tokenize() {
	log.Println("Log Obfuscator Tokenize()")
}

func (l *LogObfuscator) Parse() {
	log.Println("Log Obfuscator Parse()")
}
