package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/xuning888/ollama-hertz/internal/dal/redis"
	repo "github.com/xuning888/ollama-hertz/internal/repo/chat"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
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
	response, err = c.llm.Call(timeoutCtx, content)
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
	_, err = c.llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		sbd.Write(chunk)
		return nil
	}))
	if err != nil {
		return "", err
	}
	res = sbd.String()
	err = cache.Store(ctx, req.UserId, []*chat.Content{chat.NewContext(chat.User, req.Content), chat.NewContext(chat.Assistant, res)})
	return
}

func (c *ChatServiceImpl) ClearSession(ctx context.Context, userId string) error {
	cache := repo.NewRedisCache(redis.Client, 0)
	return cache.Clear(ctx, userId)
}

func (c *ChatServiceImpl) ChatWithSessionStream(
	ctx context.Context, req chat.ChatWithSessonReq, appCtx *app.RequestContext) (err error) {

	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(req.LlmTimeoutSecond))
	defer cancelFunc()

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

	var sbd strings.Builder
	appCtx.Response.Header.Set("mime-type", "text/event-stream")
	_, err = c.llm.GenerateContent(timeout, content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			sbd.Write(chunk)
			_, err2 := appCtx.Write(chunk)
			if err2 != nil {
				return err2
			}
			return appCtx.Flush()
		}))
	res := sbd.String()
	err = cache.Store(ctx, req.UserId, []*chat.Content{chat.NewContext(chat.User, req.Content), chat.NewContext(chat.Assistant, res)})
	return
}

func NewChatService(llm *ollama.LLM) *ChatServiceImpl {
	return &ChatServiceImpl{
		llm: llm,
	}
}
