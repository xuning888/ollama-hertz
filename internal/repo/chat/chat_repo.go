package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"time"
)

var (
	ErrorUnmarshal = errors.New("Unmarshal message content error")
	ErrorEmpty     = errors.New("Empty")
)

var _ ChatCacheRepo = (*CacheRepo)(nil)

// ChatCacheRepo
// Note: 这是一个基于redis的滑动窗口
type ChatCacheRepo interface {
	Store(ctx context.Context, userId string, contents []*chat.Content) error
	Load(ctx context.Context, userId string) (messages []*chat.Content, err error)
	Clear(ctx context.Context, userId string) error
}

type CacheRepo struct {
	client     *redis.Client
	maxWindows int
}

func (c *CacheRepo) Load(ctx context.Context, userId string) (messages []*chat.Content, err error) {
	key := c.key(userId)
	zRange := c.client.ZRange(ctx, key, 0, -1)
	if err = zRange.Err(); err != nil {
		hlog.CtxErrorf(ctx, "load chat message failed, error: %v", err)
		return
	}
	args := zRange.Val()
	if args == nil || len(args) == 0 {
		return nil, ErrorEmpty
	}
	result := make([]*chat.Content, 0, len(args))
	for _, arg := range args {
		var content chat.Content
		err := json.Unmarshal([]byte(arg), &content)
		if err != nil {
			return nil, ErrorUnmarshal
		}
		result = append(result, &content)
	}
	messages = result
	return
}

func (c *CacheRepo) Store(ctx context.Context, userId string, contents []*chat.Content) error {
	score := float64(time.Now().UnixMilli())
	key := c.key(userId)
	members := make([]redis.Z, 0, len(contents))
	for idx, content := range contents {
		message, err := json.Marshal(content)
		if err != nil {
			hlog.CtxErrorf(ctx, "store chat message failed marshal content error: %v", err)
			return err
		}
		members = append(members, redis.Z{
			Score:  score + float64(idx),
			Member: message,
		})
	}

	intCmd := c.client.ZAdd(ctx, key, members...)
	if err := intCmd.Err(); err != nil {
		hlog.CtxErrorf(ctx, "store chat message field error: %v", err)
		return err
	}
	return c.trimWindow(ctx, key)
}

func (c *CacheRepo) Clear(ctx context.Context, userId string) error {
	key := c.key(userId)
	return c.client.Del(ctx, key).Err()
}

func (c *CacheRepo) trimWindow(ctx context.Context, key string) error {
	return c.client.ZRemRangeByRank(ctx, key, 0, int64(-(c.maxWindows + 1))).Err()
}

func (c *CacheRepo) key(userId string) string {
	return fmt.Sprintf("chat:session:%s", userId)
}

func NewRedisCache(client *redis.Client, maxWindows int) *CacheRepo {
	return &CacheRepo{
		client:     client,
		maxWindows: maxWindows,
	}
}
