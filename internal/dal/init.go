package dal

import (
	"github.com/xuning888/ollama-hertz/internal/dal/mysql"
	"github.com/xuning888/ollama-hertz/internal/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
