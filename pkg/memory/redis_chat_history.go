package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"log"
	"time"
)

var (
	defaultNamespace = "chatHistory"
	defaultSessionId = "defaultSessionId"
)

type RedisMessageContent struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

func (r *RedisMessageContent) ConvertToChatMessage() llms.ChatMessage {
	switch r.Role {
	case string(llms.ChatMessageTypeAI):
		return llms.AIChatMessage{Content: r.Content}
	case string(llms.ChatMessageTypeHuman):
		return llms.HumanChatMessage{Content: r.Content}
	default:
		log.Printf("convert to chat message failed with invalid message type, type:%v\n", r.Role)
		return nil
	}
}

var _ schema.ChatMessageHistory = (*RedisChatMessageIHistory)(nil)

type RedisChatMessageIHistory struct {
	client    redis.Cmdable
	sessionId string
	namespace string
}

func (r *RedisChatMessageIHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	key := r.key()
	_message := RedisMessageContent{
		Role:      string(message.GetType()),
		Content:   message.GetContent(),
		Timestamp: time.Now().UnixMilli(),
	}
	marshalMessage, err := json.Marshal(_message)
	if err != nil {
		return err
	}
	resCmd := r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(_message.Timestamp),
		Member: marshalMessage,
	})
	if resCmd.Err() != nil {
		return resCmd.Err()
	}
	return nil
}

func (r *RedisChatMessageIHistory) AddUserMessage(ctx context.Context, message string) error {
	chatMessage := llms.HumanChatMessage{Content: message}
	return r.AddMessage(ctx, chatMessage)
}

func (r *RedisChatMessageIHistory) AddAIMessage(ctx context.Context, message string) error {
	chatMessage := llms.AIChatMessage{Content: message}
	return r.AddMessage(ctx, chatMessage)
}

func (r *RedisChatMessageIHistory) Clear(ctx context.Context) error {
	key := r.key()
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisChatMessageIHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	key := r.key()
	zRange := r.client.ZRange(ctx, key, 0, -1)
	if err := zRange.Err(); err != nil {
		return nil, err
	}
	values := zRange.Val()
	messages := make([]llms.ChatMessage, 0, len(values))
	for _, value := range values {
		var message RedisMessageContent
		if err := json.Unmarshal([]byte(value), &message); err != nil {
			return messages, err
		}
		messages = append(messages, message.ConvertToChatMessage())
	}
	return messages, nil
}

func (r *RedisChatMessageIHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	key := r.key()
	zs := make([]redis.Z, 0, len(messages))
	now := time.Now().UnixMilli()
	for idx, message := range messages {
		timestamp := int64(idx) + now
		content := RedisMessageContent{
			Role:      string(message.GetType()),
			Content:   message.GetContent(),
			Timestamp: timestamp,
		}
		member, err := json.Marshal(content)
		if err != nil {
			return err
		}
		zs = append(zs, redis.Z{
			Score:  float64(content.Timestamp),
			Member: member,
		})
	}
	pipe := r.client.TxPipeline()
	if err := pipe.Del(ctx, key).Err(); err != nil {
		return err
	}
	if err := pipe.ZAdd(ctx, key, zs...).Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RedisChatMessageIHistory) key() string {
	return fmt.Sprintf("%s:%s", r.namespace, r.sessionId)
}

func NewRedisChatMessageIHistory(opts ...RedisChatHistoryOption) (*RedisChatMessageIHistory, error) {
	rh, err := applyRedisChatHistoryOption(opts...)
	if err != nil {
		return nil, err
	}
	return rh, nil
}
