package dal

import (
	"github.com/xuning888/yoyoyo/internal/dal/database"
	"github.com/xuning888/yoyoyo/internal/dal/redis"
)

func Init() {
	redis.Init()
	database.Init()
}
