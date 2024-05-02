package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xuning888/yoyoyo/internal/schema/prompt"
	"github.com/xuning888/yoyoyo/internal/service"
	"github.com/xuning888/yoyoyo/pkg/api"
	"github.com/xuning888/yoyoyo/pkg/logger"
	"net/http"
)

type PromptController struct {
	promptService service.PromptService
	lg            logger.Logger
}

func (p *PromptController) PromptPageInfo(c *gin.Context) {
	defer p.lg.Sync()

	var request prompt.PromptPageReq

	if err := c.BindQuery(&request); err != nil {
		p.lg.Errorf("PromptPageInfo request invalid json error: %v", err)
		c.JSON(http.StatusOK, api.FailedWithMessage("Invalid JSON"))
		return
	}

	prompts, total, err := p.promptService.PromptPage(c, &request)
	if err != nil {
		p.lg.Errorf("PromptPageInfo prompt request failed with error: %v", err)
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
	defer p.lg.Sync()
	var request prompt.PromptAddOrUpdate

	if err := c.BindJSON(&request); err != nil {
		p.lg.Errorf("PromptAddOrUpdate invalid json error: %v", err)
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
		lg:            logger.Named("PromptController"),
	}
}
