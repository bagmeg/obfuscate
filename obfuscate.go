package obfuscate

import (
	"github.com/bagmeg/obfuscate/factory/log"
	"github.com/bagmeg/obfuscate/factory/sql"
)

func NewObfuscator(t string) Obfuscator {
	switch t {
	case "db":
		return sql.NewSQLObfuscator()
	case "log":
		return log.NewLogObfuscator()
	default:
		return nil
	}
}
