package controller

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gin-gonic/gin"
	"github.com/xuning888/ollama-hertz/internal/schema/prompt"
	"github.com/xuning888/ollama-hertz/internal/service"
	"github.com/xuning888/ollama-hertz/pkg/api"
	"net/http"
)

type PromptController struct {
	promptService service.PromptService
}

func (p *PromptController) PromptPageInfo(c *gin.Context) {

	var request prompt.PromptPageReq

	if err := c.BindQuery(&request); err != nil {
		hlog.CtxErrorf(c, "chat fialed invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	prompts, total, err := p.promptService.PromptPage(c, &request)
	if err != nil {
		hlog.CtxErrorf(c, "prompt request faield with error: %v", err)
		c.JSON(http.StatusOK, api.Failed())
		return
	}

	pageInfo := &api.PageInfo{
		PageNum:  request.PageNum,
		PageSize: request.PageSize,
		Total:    total,
		List:     prompts,
	}
	c.JSON(http.StatusOK, api.SuccessWithData(pageInfo))
}

func (p *PromptController) PromptAddOrUpdate(c *gin.Context) {
	var request prompt.PromptAddOrUpdate

	if err := c.BindJSON(&request); err != nil {
		hlog.CtxErrorf(c, "promptAddOrUpdate Invalid JSON error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}

	err := p.promptService.AddOrUpdate(c, &request)
	if err != nil {
		c.JSON(http.StatusOK, api.FailedWithMessage(err.Error()))
		return
	}
	c.JSON(http.StatusOK, api.Success())
	return
}

func (p *PromptController) Register(server *gin.Engine) {
	promptGroup := server.Group("/api/v1/prompt")
	promptGroup.GET("/page", p.PromptPageInfo)
	promptGroup.POST("/addOrUpdate", p.PromptAddOrUpdate)
}

func NewPromptController(promptService service.PromptService) *PromptController {
	return &PromptController{
		promptService: promptService,
	}
}
