package service

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/xuning888/ollama-hertz/internal/dal/redis"
	repo "github.com/xuning888/ollama-hertz/internal/repo/chat"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"github.com/xuning888/ollama-hertz/pkg/api"
	"strings"
	"time"
)

var _ ChatService = (*ChatServiceImpl)(nil)

var (
	ErrorCallLlmTimeout = errors.New("call llm timeout")
)

type ChatService interface {
	Chat(ctx context.Context, req chat.ChatReq) (response string, err error)
	ChatWithSession(ctx context.Context, req chat.ChatWithSessonReq) (response string, err error)
	ChatWithSessionStream(ctx context.Context, req chat.ChatWithSessonReq, appCtx *app.RequestContext) (err error)
	ClearSession(ctx context.Context, userId string) error
}

type ChatServiceImpl struct {
	llm *ollama.LLM
}

func (c *ChatServiceImpl) Chat(ctx context.Context, req chat.ChatReq) (response string, err error) {
	timeoutSecond := req.LlmTimeoutSecond
	content := req.Content
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(timeoutSecond))
	defer cancelFunc()

	response, err = llms.GenerateFromSinglePrompt(timeoutCtx, c.llm, req.Content)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			hlog.CtxErrorf(ctx, "chat faield, call llm timeout content: %v, error: %v", err, err)
			err = ErrorCallLlmTimeout
		} else {
			hlog.CtxErrorf(ctx, "chat failed contet: %v with error: %v", content, err)
		}
		return
	}
	return
}

func (c *ChatServiceImpl) ChatWithSession(ctx context.Context, req chat.ChatWithSessonReq) (res string, err error) {

	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(req.LlmTimeoutSecond))
	defer cancelFunc()

	// user 的时间戳
	userMillis := time.Now().UnixMilli()

	cache := repo.NewRedisCache(redis.Client, req.MaxWindows)

	messages, err := cache.Load(ctx, req.UserId)
	if err != nil {
		if !errors.Is(err, repo.ErrorEmpty) {
			return "", err
		}
	}

	content := make([]llms.MessageContent, 0, len(messages))
	for _, msg := range messages {
		content = append(content, llms.TextParts(msg.Role.LlmsRole(), msg.Message))
	}
	content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, req.Content))

	var sbd strings.Builder
	_, err = c.llm.GenerateContent(timeout, content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			sbd.Write(chunk)
			return nil
		}))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", ErrorCallLlmTimeout
		}
		return "", err
	}
	assistantMillis := time.Now().UnixMilli()
	res = sbd.String()
	err = cache.Store(ctx, req.UserId, []*chat.Content{
		chat.NewContext(chat.User, req.Content, userMillis),
		chat.NewContext(chat.Assistant, res, assistantMillis)})
	return
}

func (c *ChatServiceImpl) ClearSession(ctx context.Context, userId string) error {
	cache := repo.NewRedisCache(redis.Client, 0)
	return cache.Clear(ctx, userId)
}

func (c *ChatServiceImpl) ChatWithSessionStream(
	ctx context.Context, req chat.ChatWithSessonReq, appCtx *app.RequestContext) (err error) {

	userMillis := time.Now().UnixMilli()

	cache := repo.NewRedisCache(redis.Client, req.MaxWindows)

	messages, err := cache.Load(ctx, req.UserId)
	if err != nil {
		if !errors.Is(err, repo.ErrorEmpty) {
			return err
		}
	}

	content := make([]llms.MessageContent, 0, len(messages))
	for _, msg := range messages {
		content = append(content, llms.TextParts(msg.Role.LlmsRole(), msg.Message))
	}
	content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, req.Content))

	appCtx.Response.Header.Set("Content-Type", "application/x-ndjson")
	appCtx.Response.Header.Set("Connection", "Keep-Alive")
	appCtx.Response.Header.Set("X-Content-Type-Options", "nosniff")

	resp, err := c.llm.GenerateContent(ctx, content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			responseData := api.SuccessWithData(string(chunk))
			bytes, _ := json.Marshal(responseData)
			_, err2 := appCtx.Write(bytes)
			if err2 != nil {
				return err2
			}
			_, err3 := appCtx.Write([]byte{'\n'})
			if err3 != nil {
				return err3
			}
			return appCtx.Flush()
		}),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrorCallLlmTimeout
		}
		hlog.CtxErrorf(ctx, "call llm has error: %v", err)
		return err
	}
	assistantMillis := time.Now().UnixMilli()
	res := resp.Choices[0].Content
	err = cache.Store(ctx, req.UserId, []*chat.Content{
		chat.NewContext(chat.User, req.Content, userMillis),
		chat.NewContext(chat.Assistant, res, assistantMillis)})
	return
}

func (c *ChatServiceImpl) Embedding(ctx context.Context) {
	c.llm.CreateEmbedding(ctx, nil)
}

func NewChatService(llm *ollama.LLM) *ChatServiceImpl {
	return &ChatServiceImpl{
		llm: llm,
	}
}
