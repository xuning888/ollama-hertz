package service

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/xuning888/ollama-hertz/internal/dal/redis"
	repo "github.com/xuning888/ollama-hertz/internal/repo/chat"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"github.com/xuning888/ollama-hertz/pkg/api"
	"github.com/xuning888/ollama-hertz/pkg/config"
	"strings"
	"time"
)

var _ ChatService = (*ChatServiceImpl)(nil)

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
}

func (c *ChatServiceImpl) ClearSession(ctx context.Context, userId string) error {
	cache := repo.NewRedisCache(redis.Client, 0)
	return cache.Clear(ctx, userId)
}

func (c *ChatServiceImpl) ChatWithSessionStream(
	ctx context.Context, req chat.ChatWithSessonReq, appCtx *gin.Context) (err error) {

	userId, llmModel, maxWindows := req.UserId, req.LlmModel, req.MaxWindows

	llm, err := ollama.New(ollama.WithModel(llmModel),
		ollama.WithServerURL(config.DefaultConfig.OllmServerUrl))
	if err != nil {
		hlog.CtxErrorf(ctx, "create ollama llm error: %v", err)
		return err
	}

	userMillis := time.Now().UnixMilli()
	cache := repo.NewRedisCache(redis.Client, maxWindows)

	// 获取对话窗口
	messages, err := cache.Load(ctx, userId)
	if err != nil {
		if !errors.Is(err, repo.ErrorEmpty) {
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
		llms.WithStreamingFunc(makeStreamHandler(appCtx)),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrorCallLlmTimeout
		}
		hlog.CtxErrorf(ctx, "call llm has error: %v", err)
		return errors.New("call llm failed")
	}
	assistantMillis := time.Now().UnixMilli()
	res := resp.Choices[0].Content
	err = cache.Store(ctx, req.UserId, []*chat.Content{
		chat.NewContext(chat.User, req.Content, userMillis),
		chat.NewContext(chat.Assistant, res, assistantMillis)})
	if err != nil {
		return err
	}
	return
}

func (c *ChatServiceImpl) summary(ctx context.Context, content []llms.MessageContent) (summary string, err error) {

	content = append(content, llms.TextParts(llms.ChatMessageTypeSystem,
		"为了帮助AI理解对话内容，为这一些列对话内容生成摘要"))
	var sbd strings.Builder
	llm, err := ollama.New(
		ollama.WithModel(config.DefaultConfig.LlmModel),
		ollama.WithServerURL(config.DefaultConfig.OllmServerUrl),
		ollama.WithRunnerMainGPU(8),
	)
	_, err = llm.GenerateContent(context.Background(), content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			sbd.Write(chunk)
			return nil
		}))
	if err != nil {
		return "", err
	}
	summary = sbd.String()
	hlog.CtxInfof(ctx, "生成摘要内容: %v", summary)
	return
}

func makeStreamHandler(c *gin.Context) func(ctx context.Context, chunk []byte) error {
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Connection", "Keep-Alive")
	c.Header("X-Content-Type-Options", "nosniff")
	return func(ctx context.Context, chunk []byte) error {
		responseData := api.SuccessWithData(string(chunk))
		bytes, _ := json.Marshal(responseData)
		_, err := c.Writer.Write(bytes)
		if err != nil {
			return err
		}
		_, err = c.Writer.Write([]byte{'\n'})
		if err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
}

func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{}
}
