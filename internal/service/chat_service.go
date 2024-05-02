package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/xuning888/yoyoyo/config"
	"github.com/xuning888/yoyoyo/internal/dal/redis"
	repo "github.com/xuning888/yoyoyo/internal/repo/chat"
	"github.com/xuning888/yoyoyo/internal/schema/chat"
	"github.com/xuning888/yoyoyo/pkg/api"
	"github.com/xuning888/yoyoyo/pkg/logger"
	"net/http"
	"time"
)

var _ ChatService = (*ChatServiceImpl)(nil)

var client = &http.Client{
	Transport: http.DefaultTransport,
}

var (
	ErrorCallLlmTimeout = errors.New("call llm timeout")
)

type ChatService interface {
	// ChatWithSessionStream QA with stream
	ChatWithSessionStream(ctx context.Context, req chat.ChatWithSessonReq, appCtx *gin.Context) (err error)
	// ClearSession clear chat message context
	ClearSession(ctx context.Context, userId string) error
}

type ChatServiceImpl struct {
	lg logger.Logger
}

func (c *ChatServiceImpl) ClearSession(ctx context.Context, userId string) error {
	cache := repo.NewRedisCache(redis.Client, 0)
	return cache.Clear(ctx, userId)
}

func (c *ChatServiceImpl) ChatWithSessionStream(
	ctx context.Context, req chat.ChatWithSessonReq, appCtx *gin.Context) (err error) {
	userId, llmModel, maxWindows := req.UserId, req.LlmModel, req.MaxWindows
	llm, err := ollama.New(
		ollama.WithModel(llmModel),
		ollama.WithHTTPClient(client),
		ollama.WithServerURL(config.DefaultConfig.OllmServerUrl),
	)
	if err != nil {
		c.lg.Errorf("ChatWithSessionStream create ollama llm error: %v", err)
		return err
	}
	userMillis := time.Now().UnixMilli()
	cache := repo.NewRedisCache(redis.Client, maxWindows)

	// 获取对话窗口
	messages, err := cache.Load(ctx, userId)
	if err != nil {
		if !errors.Is(err, repo.ErrorEmpty) {
			c.lg.Errorf("ChatWithSessionStream query message windows error: %v", err)
			return err
		}
	}
	// 历史消息
	historyMessage := make([]llms.MessageContent, 0, len(messages))
	for _, msg := range messages {
		historyMessage = append(historyMessage, llms.TextParts(msg.Role.LlmsRole(), msg.Message))
	}
	// 当前对话内容
	curCtxMessage := append(historyMessage, llms.TextParts(llms.ChatMessageTypeHuman, req.Content))
	// call llm
	resp, err := llm.GenerateContent(ctx, curCtxMessage,
		llms.WithStreamingFunc(c.makeStreamHandler(appCtx)),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.lg.Errorf("call llm failed with error: %v", err)
			return ErrorCallLlmTimeout
		}
		c.lg.Errorf("call llm failed with error: %v", err)
		return errors.New("call llm failed")
	}
	assistantMillis := time.Now().UnixMilli()
	choice := resp.Choices[0]
	info := choice.GenerationInfo
	completionTokens, promptTokens, totalTokens := info["CompletionTokens"], info["PromptTokens"], info["TotalTokens"]
	c.lg.Infof("totalTokens: %v, PromptTokens: %v, completionTokens: %v", totalTokens, promptTokens, completionTokens)

	res := choice.Content
	err = cache.Store(ctx, req.UserId, []*chat.Content{
		chat.NewContext(chat.User, req.Content, userMillis),
		chat.NewContext(chat.Assistant, res, assistantMillis)})
	if err != nil {
		return err
	}
	return
}

func (c *ChatServiceImpl) makeStreamHandler(appCtx *gin.Context) func(ctx context.Context, chunk []byte) error {
	appCtx.Header("Content-Type", "application/x-ndjson")
	appCtx.Header("Connection", "Keep-Alive")
	appCtx.Header("X-Content-Type-Options", "nosniff")
	return func(ctx context.Context, chunk []byte) error {
		responseData := api.SuccessWithData(string(chunk))
		bytes, _ := json.Marshal(responseData)
		_, err1 := appCtx.Writer.Write(bytes)
		if err1 != nil {
			return fmt.Errorf("write stream response error: %w", err1)
		}
		_, err2 := appCtx.Writer.Write([]byte{'\n'})
		if err2 != nil {
			return fmt.Errorf("write stream LF error: %w", err2)
		}
		appCtx.Writer.Flush()
		return nil
	}
}

func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{
		lg: logger.Named("ChatServiceImpl"),
	}
}
