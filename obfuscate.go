package obfuscate

import (
	"github.com/bagmeg/obfuscate/config"
	"github.com/bagmeg/obfuscate/factory/log"
	"github.com/bagmeg/obfuscate/factory/sql"
)

func NewObfuscator(cfg *config.Config) (Obfuscator, Obfuscator) {
	switch {
	case cfg.Enabled.Sql:
		if cfg.Enabled.Log {
			return sql.NewSQLObfuscator(), log.NewLogObfuscator()
		}
		return sql.NewSQLObfuscator(), nil
	case cfg.Enabled.Log:
		return nil, log.NewLogObfuscator()
	default:
		return nil, nil
	}
}
