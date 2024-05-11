package memory

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrorNilRedisCmdable = errors.New("redis cmdable is empty")
)

type RedisChatHistoryOption func(history *RedisChatMessageIHistory)

func applyRedisChatHistoryOption(opts ...RedisChatHistoryOption) (*RedisChatMessageIHistory, error) {
	rh := &RedisChatMessageIHistory{}
	for _, opt := range opts {
		opt(rh)
	}
	if rh.client == nil {
		return nil, ErrorNilRedisCmdable
	}

	// check redis client
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	if err := rh.client.Ping(timeout).Err(); err != nil {
		return nil, err
	}

	if rh.namespace == "" {
		rh.namespace = defaultNamespace
	}

	if rh.sessionId == "" {
		rh.sessionId = defaultSessionId
	}
	return rh, nil
}

func WithRedisCmdable(cmdable redis.Cmdable) RedisChatHistoryOption {
	return func(history *RedisChatMessageIHistory) {
		history.client = cmdable
	}
}

func WithSessionId(sessionId string) RedisChatHistoryOption {
	return func(history *RedisChatMessageIHistory) {
		history.sessionId = sessionId
	}
}

func WithNamespace(namespace string) RedisChatHistoryOption {
	return func(history *RedisChatMessageIHistory) {
		history.namespace = namespace
	}
}
