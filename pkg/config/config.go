package config

import "gorm.io/gorm/logger"

var DefaultConfig = &Config{}

type Config struct {
	LlmModel   string
	ServerPort string
	Redis      RedisCfg
	MySQL      MySQLCfg
}

type RedisCfg struct {
	Addr     string
	Password string
	DB       int
}

type MySQLCfg struct {
	DSN      string
	Username string
	Password string
	LogModel logger.LogLevel
}
