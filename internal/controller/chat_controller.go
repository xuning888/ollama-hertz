package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xuning888/ollama-hertz/internal/schema/chat"
	"github.com/xuning888/ollama-hertz/internal/service"
	"github.com/xuning888/ollama-hertz/pkg/api"
	"github.com/xuning888/ollama-hertz/pkg/logger"
	"net/http"
)

type ChatController struct {
	chatService service.ChatService
	lg          logger.Logger
}

func (cc *ChatController) ChatSessionStream(c *gin.Context) {
	var request chat.ChatWithSessonReq
	if err := c.BindJSON(&request); err != nil {
		cc.lg.Errorf("ChatSessionStream failed invalid json error:%v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	err := cc.chatService.ChatWithSessionStream(c, request, c)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	return
}

func (cc *ChatController) ChatClearSession(c *gin.Context) {
	var request struct {
		UserId string
	}

	if err := c.BindJSON(&request); err != nil {
		cc.lg.Errorf("ChatClearSession failed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	cc.lg.Infof("ChatClearSession clear chat context, request: %v", request)

	err := cc.chatService.ClearSession(c, request.UserId)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage("上下文清理失败"))
		return
	}
	c.JSON(http.StatusOK, api.Success())
}

func (cc *ChatController) Register(engine *gin.Engine) {
	apiV1 := engine.Group("/api/v1/chat")
	apiV1.POST("/stream", cc.ChatSessionStream)
	apiV1.POST("/stream/clear", cc.ChatClearSession)
}

func NewChatController(chatService service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
		lg:          logger.Named("ChatController"),
	}
}
