package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/xuning888/yoyoyo/config"
)

var Client redis.Cmdable

type Model string

var (
	standalone Model = "standalone"
	cluster    Model = "cluster"
	sentinel   Model = "sentinel"
)

func Init() {
	redisCfg := config.DefaultConfig.Redis
	switch redisCfg.Model {
	case string(standalone):
		if redisCfg.Addr == "" {
			panic(fmt.Sprintf("standalone redis addr is emtpy"))
		}
		Client = redis.NewClient(&redis.Options{
			Addr:     config.DefaultConfig.Redis.Addr,
			Password: config.DefaultConfig.Redis.Password,
			DB:       config.DefaultConfig.Redis.DB,
		})
		break
	case string(cluster):
		if redisCfg.Addrs == nil || len(redisCfg.Addrs) == 0 {
			panic(fmt.Sprintf("cluster redis addrs is emtpy"))
		}
		Client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.Addrs,
			Password: redisCfg.Password,
		})
		break
	case string(sentinel):
		if redisCfg.SentinelAddrs == nil || len(redisCfg.SentinelAddrs) == 0 {
			panic(fmt.Sprintf("sentinel addrs is empty"))
		}
		if redisCfg.MasterName == "" {
			panic(fmt.Sprintf("sentinel masterName is empty"))
		}
		redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    redisCfg.MasterName,
			SentinelAddrs: redisCfg.SentinelAddrs,
		})
		break
	default:
		panic(fmt.Sprintf("unknow redis model: %v", redisCfg.Model))
	}
}
