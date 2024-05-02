package database

import (
	"github.com/xuning888/ollama-hertz/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.DefaultConfig.MySQL.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(config.DefaultConfig.MySQL.LogModel),
	})
	if err != nil {
		panic(err)
	}
}
