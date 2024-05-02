package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"github.com/xuning888/ollama-hertz/pkg/logger"
)

var (
	ErrorUnmarshal = errors.New("Unmarshal message content error")
	ErrorEmpty     = errors.New("Empty")
)

var _ Cache = (*RedisCache)(nil)

// Cache
// Note: 这是一个基于redis的滑动窗口
type Cache interface {
	Store(ctx context.Context, userId string, contents []*chat.Content) error
	Load(ctx context.Context, userId string) (messages []*chat.Content, err error)
	Clear(ctx context.Context, userId string) error
}

type RedisCache struct {
	client     redis.Cmdable
	maxWindows int
	lg         logger.Logger
}

func (c *RedisCache) Load(ctx context.Context, userId string) (messages []*chat.Content, err error) {
	defer c.lg.Sync()
	key := c.key(userId)
	c.lg.Infof("load history messages key: %s", key)
	zRange := c.client.ZRange(ctx, key, 0, -1)
	if err = zRange.Err(); err != nil {
		c.lg.Errorf("load chat history messages failed, key: %s error: %v", key, err)
		return
	}
	values := zRange.Val()
	if values == nil || len(values) == 0 {
		c.lg.Infof("load chat history messages empty, key: %s", key)
		return nil, ErrorEmpty
	}
	result := make([]*chat.Content, 0, len(values))
	for _, arg := range values {
		var content chat.Content
		err := json.Unmarshal([]byte(arg), &content)
		if err != nil {
			c.lg.Errorf("Load messages json unmarshal key: %s, error: %v", key, err)
			return nil, ErrorUnmarshal
		}
		result = append(result, &content)
	}
	messages = result
	c.lg.Infof("load chat history messages success size: %d, key: %s", len(messages), key)
	return
}

func (c *RedisCache) Store(ctx context.Context, userId string, contents []*chat.Content) (err error) {
	defer c.lg.Sync()
	key := c.key(userId)
	if len(contents) == 0 {
		c.lg.Infof("store chat messages is empty key: %s", key)
		return
	}
	c.lg.Infof("sotre chat messages key: %s, messages size: %d", key, len(contents))
	members := make([]redis.Z, 0, len(contents))
	for _, content := range contents {
		message, err2 := json.Marshal(content)
		if err2 != nil {
			c.lg.Errorf("Store chat messages failed marshal key: %s, content: %v error: %v", key, content, err)
			return err2
		}
		members = append(members, redis.Z{Score: float64(content.Timestamp), Member: message})
	}
	intCmd := c.client.ZAdd(ctx, key, members...)
	if err3 := intCmd.Err(); err3 != nil {
		c.lg.Errorf("store chat messages field key: %s error: %v", key, err)
		return err3
	}
	return c.trimWindow(ctx, key)
}

func (c *RedisCache) Clear(ctx context.Context, userId string) (err error) {
	defer c.lg.Sync()
	key := c.key(userId)
	c.lg.Infof("clear session key: %v", key)
	err = c.client.Del(ctx, key).Err()
	if err != nil {
		c.lg.Errorf("clear session key: %v error: %v", key, err)
		return err
	}
	return
}

func (c *RedisCache) trimWindow(ctx context.Context, key string) error {
	return c.client.ZRemRangeByRank(ctx, key, 0, int64(-(c.maxWindows + 1))).Err()
}

func (c *RedisCache) key(userId string) string {
	return fmt.Sprintf("chat:session:%s", userId)
}

func NewRedisCache(client redis.Cmdable, maxWindows int) *RedisCache {
	return &RedisCache{
		client:     client,
		maxWindows: maxWindows,
		lg:         logger.Named("RedisCache"),
	}
}
