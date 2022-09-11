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
			return sql.NewSQLObfuscator(cfg.SQLConfig), log.NewLogObfuscator(cfg.LogConfig)
		}
		return sql.NewSQLObfuscator(cfg.SQLConfig), nil
	case cfg.Enabled.Log:
		return nil, log.NewLogObfuscator(cfg.LogConfig)
	default:
		return nil, nil
	}
}
