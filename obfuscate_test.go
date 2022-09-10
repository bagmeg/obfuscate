package obfuscate

import (
	"github.com/bagmeg/obfuscate/config"
	"log"
	"testing"
)

func TestObfuscate(t *testing.T) {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		t.Errorf("Failed to load %v", "./config.yaml")
	}

	SQLobfuscator, LogObfuscator := NewObfuscator(&cfg)

	if LogObfuscator == nil {
		log.Println("Log Obfuscator is disabled")
	}
	if SQLobfuscator == nil {
		log.Println("SQL Obfuscator is diabled")
	}

	if SQLobfuscator != nil {
		SQLobfuscator.Tokenize("query here")
		SQLobfuscator.Parse()
	}
	if LogObfuscator != nil {
		LogObfuscator.Tokenize("some log")
		LogObfuscator.Parse()
	}
}
