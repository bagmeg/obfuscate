package log

import (
	"github.com/bagmeg/obfuscate/config"
	"log"
)

type LogObfuscator struct {
	tokenizer LogTokeniezr
	cfg       *config.LogConfig
}

func NewLogObfuscator(cfg config.LogConfig) *LogObfuscator {
	return &LogObfuscator{
		cfg: &cfg,
	}
}

func (l *LogObfuscator) Tokenize(log string) {
	l.tokenizer.Tokenize()
}

func (l *LogObfuscator) Parse() {
	log.Println("Log Obfuscator Parse()")
}
