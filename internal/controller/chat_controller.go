package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xuning888/yoyoyo/internal/schema/chat"
	"github.com/xuning888/yoyoyo/internal/service"
	"github.com/xuning888/yoyoyo/pkg/api"
	"github.com/xuning888/yoyoyo/pkg/logger"
	"net/http"
)

type ChatController struct {
	chatService service.ChatService
	lg          logger.Logger
}

func (cc *ChatController) ChatSessionStream(c *gin.Context) {
	defer cc.lg.Sync()
	var request chat.ChatWithSessonReq
	if err := c.BindJSON(&request); err != nil {
		cc.lg.Errorf("ChatSessionStream failed invalid json error:%v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	err := cc.chatService.ChatWithSessionStreamV1(c, request, c)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	return
}

func (cc *ChatController) UploadFile(c *gin.Context) {
	// 200MB
	const MaxUploadFileSize = 200 << 20
	// 限制文件读取的大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadFileSize)

	if err := c.Request.ParseMultipartForm(MaxUploadFileSize); err != nil {
		cc.lg.Errorf("UploadFile error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("文件太大了"))
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		cc.lg.Errorf("获取文件失败, error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("无效请求"+err.Error()))
		return
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			cc.lg.Errorf("关闭文件流失败, error: %v", closeErr)
		}
	}()

	if err2 := c.SaveUploadedFile(fileHeader, "./uploads/"+fileHeader.Filename); err2 != nil {
		cc.lg.Errorf("save file error: %v", err2)
		c.JSON(http.StatusOK, api.FailedWithMessage("upload file error: "+err2.Error()))
		return
	}

	c.String(http.StatusOK, "文件上传成功！")
}

func (cc *ChatController) ChatClearSession(c *gin.Context) {
	defer cc.lg.Sync()
	var request struct {
		ChatId string
		UserId string
	}

	if err := c.BindJSON(&request); err != nil {
		cc.lg.Errorf("ChatClearSession failed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	cc.lg.Infof("ChatClearSession clear chat context, request: %v", request)

	err := cc.chatService.ClearSession(c, request.ChatId, request.UserId)
	if err != nil {
		cc.lg.Errorf("clear session error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("clear session error"))
		return
	}
	c.JSON(http.StatusOK, api.Success())
}

func (cc *ChatController) Register(engine *gin.Engine) {
	apiV1 := engine.Group("/api/v1/chat")
	apiV1.POST("/stream", cc.ChatSessionStream)
	apiV1.POST("/stream/clear", cc.ChatClearSession)
	apiV1.POST("/upload", cc.UploadFile)
}

func NewChatController(chatService service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
		lg:          logger.Named("ChatController"),
	}
}
