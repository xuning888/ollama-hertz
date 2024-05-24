package lock

import (
	"github.com/go-redsync/redsync/v4"
	gredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

var Rs *redsync.Redsync

func Init(client redis.UniversalClient) {
	pool := gredis.NewPool(client)
	Rs = redsync.New(pool)
}
