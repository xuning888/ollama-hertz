package dal

import (
	"github.com/xuning888/ollama-hertz/internal/dal/database"
	"github.com/xuning888/ollama-hertz/internal/dal/redis"
)

func Init() {
	redis.Init()
	database.Init()
}
