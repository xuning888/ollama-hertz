package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"github.com/xuning888/ollama-hertz/internal/service"
	"github.com/xuning888/ollama-hertz/pkg/api"
	"net/http"
)

type ChatController struct {
	chatService service.ChatService
}

func (cc *ChatController) Chat(ctx context.Context, c *app.RequestContext) {
	var request chat.ChatReq

	if err := c.BindJSON(&request); err != nil {
		hlog.CtxErrorf(ctx, "chat fialed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}
	response, err := cc.chatService.Chat(ctx, request)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	c.JSON(http.StatusOK, api.SuccessWithData(response))
	return
}

func (cc *ChatController) ChatWithSession(ctx context.Context, c *app.RequestContext) {
	var request chat.ChatWithSessonReq

	if err := c.BindJSON(&request); err != nil {
		hlog.CtxErrorf(ctx, "chatWithSession fialed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	response, err := cc.chatService.ChatWithSession(ctx, request)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	c.JSON(http.StatusOK, api.SuccessWithData(response))
	return
}

func (cc *ChatController) ChatStream(ctx context.Context, c *app.RequestContext) {
	var request chat.ChatWithSessonReq
	if err := c.BindJSON(&request); err != nil {
		hlog.CtxErrorf(ctx, "chatWithSession fialed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}
	err := cc.chatService.ChatWithSessionStream(ctx, request, c)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	return
}

func (cc *ChatController) ChatPdf(ctx context.Context, c *app.RequestContext) {

}

func (cc *ChatController) ChatClearSession(ctx context.Context, c *app.RequestContext) {
	var request struct {
		UserId string
	}

	if err := c.BindJSON(&request); err != nil {
		hlog.CtxErrorf(ctx, "chatWithSession fialed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	hlog.CtxInfof(ctx, "clear chat context, request: %v", request)

	err := cc.chatService.ClearSession(ctx, request.UserId)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage("上下文清理失败"))
		return
	}
	c.JSON(http.StatusOK, api.Success())
}

func (cc *ChatController) Register(hertz *server.Hertz) {
	apiV1 := hertz.Group("/api/v1/chat")
	apiV1.POST("/", cc.Chat)
	apiV1.POST("/session", cc.ChatWithSession)
	apiV1.POST("/stream", cc.ChatStream)
	apiV1.POST("/stream/clear", cc.ChatClearSession)
}

func NewChatController(chatService service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}
