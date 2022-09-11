package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type DBMSKIND int

const (
	DBMSERROR      = iota
	DBMSMARIA      = iota
	DBMSMYSQL      = iota
	DBMSPOSTGRESQL = iota

	MARIA      = "MARIA"
	MYSQL      = "MYSQL"
	POSTGRESQL = "POSTGRESQL"
)

var dbms = map[string]DBMSKIND{
	MARIA:      DBMSMARIA,
	MYSQL:      DBMSMYSQL,
	POSTGRESQL: DBMSPOSTGRESQL,
}

// Load loads obfuscate configuration.
// replace digit is enabled by default
func Load(path string) (Config, error) {
	var err error
	cfg := Config{
		SQLConfig: SQLConfig{ReplaceDigits: true},
		LogConfig: LogConfig{ReplaceDigits: true},
	}

	yamlBytes, err := ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		return Config{}, err
	}
	err = cfg.setValue()
	return cfg, nil
}

type Config struct {
	Enabled   List      `yaml:"enable_list"`
	SQLConfig SQLConfig `yaml:"sql_config,omitempty"`
	LogConfig LogConfig `yaml:"log_config,omitempty"`
}

type List struct {
	Sql bool `yaml:"sql,omitempty"`
	Log bool `yaml:"log,omitempty"`
}

type LogConfig struct {
	Enabled       bool
	ReplaceDigits bool `yaml:"digit,omitempty"`
}

type SQLConfig struct {
	Enabled       bool
	DBMS          DBMSKIND
	DbKind        string `yaml:"db_kind"`
	ReplaceDigits bool   `yaml:"digit,omitempty"`
}

func (cfg *Config) setValue() error {
	if err := cfg.SQLConfig.setValue(cfg.Enabled.Sql); err != nil {
		return err
	}
	if err := cfg.LogConfig.setValue(cfg.Enabled.Log); err != nil {
		return err
	}
	return nil
}

// setValue sets database kind based on user provided kind
// if user provided kind is not mysql, postgresql, maria then error is returned
func (s *SQLConfig) setValue(enable bool) error {
	dbKind := strings.ToUpper(s.DbKind)
	val, ok := dbms[dbKind]
	if ok {
		s.DBMS = val
	} else {
		s.DBMS = DBMSERROR
		return fmt.Errorf("DBMS %s not supported", dbKind)
	}
	s.Enabled = enable
	return nil
}

func (l *LogConfig) setValue(enable bool) error {
	l.Enabled = enable
	return nil
}
