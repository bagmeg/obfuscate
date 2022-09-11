package log

import (
	"github.com/bagmeg/obfuscate/config"
)

type LogObfuscator struct {
	scanner LogScanner
	cfg     *config.LogConfig
}

func NewLogObfuscator(cfg config.LogConfig) *LogObfuscator {
	return &LogObfuscator{
		cfg: &cfg,
	}
}

func (l *LogObfuscator) Scan(log string) (string, error) {
	l.scanner.Scan()
	return "", nil
}
