package config

import "gorm.io/gorm/logger"

var DefaultConfig = &Config{}

type Config struct {
	LlmModel      string
	OllmServerUrl string
	ServerPort    string
	Redis         RedisCfg
	MySQL         MySQLCfg
}

type RedisCfg struct {
	// Model standalone(单机模式)  cluster(集群模式) sentinel(哨兵模式)
	Model string
	// Addr 如果redis的模式是单机模式, 那么就设置这个字段
	Addr string
	// Addrs 如果redis的模式是集群模式，就设置这个字段
	Addrs []string
	// SentinelAddrs 如果redis的模式是哨兵模式，就需要设置哨兵的地址
	SentinelAddrs []string
	// MasterName 如果redis的模式是哨兵模式，就需要设置 master
	MasterName string
	Password   string
	DB         int
}

type MySQLCfg struct {
	DSN      string
	LogModel logger.LogLevel
}
