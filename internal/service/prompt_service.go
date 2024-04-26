package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xuning888/ollama-hertz/internal/dal/mysql"
	modelPrompt "github.com/xuning888/ollama-hertz/internal/model/prompt"
	"github.com/xuning888/ollama-hertz/internal/schema/prompt"
)

var (
	_ PromptService = (*PromptServiceImpl)(nil)
)

type PromptService interface {

	// PromptPage
	// Note: 分页查询prompts
	PromptPage(ctx context.Context, req *prompt.PromptPageReq) (prompts []*prompt.PromptPageRes, total int64, err error)

	// AddOrUpdate
	// Note: 新增或修改Prompt
	AddOrUpdate(ctx context.Context, req *prompt.PromptAddOrUpdate) error
}

type PromptServiceImpl struct{}

func (p *PromptServiceImpl) AddOrUpdate(ctx context.Context, req *prompt.PromptAddOrUpdate) error {
	entity := &modelPrompt.Prompt{}
	entity.ID = req.Id
	entity.Name = req.Name
	entity.Description = req.Description
	entity.Temperature = req.Temperature
	if entity.ID == uint(0) {
		// add
		err := mysql.PromptAdd(ctx, entity)
		if err != nil {
			hlog.CtxErrorf(ctx, "add prompt failed with error: %v", err)
			return err
		}
		return nil
	}

	err := mysql.PromptUpdateById(ctx, entity)
	if err != nil {
		hlog.CtxErrorf(ctx, "update prompt fialed with error: %v", err)
		return err
	}
	return nil
}

func (p *PromptServiceImpl) PromptPage(ctx context.Context, req *prompt.PromptPageReq) (
	prompts []*prompt.PromptPageRes, total int64, err error) {

	pageSize, pageNum, name := req.PageSize, req.PageNum, req.Name

	dbPrompts, total, err := mysql.PromptPageInfo(ctx, pageSize, pageNum, name)
	if err != nil {
		return
	}

	result := make([]*prompt.PromptPageRes, 0, len(dbPrompts))

	for _, dbPrompt := range dbPrompts {
		res := &prompt.PromptPageRes{}
		res.Id = dbPrompt.ID
		res.Name = dbPrompt.Name
		res.Description = dbPrompt.Description
		res.Temperature = dbPrompt.Temperature
		res.CreatedAt = dbPrompt.CreatedAt
		res.UpdatedAt = dbPrompt.UpdatedAt
		result = append(result, res)
	}
	prompts = result
	return
}

func NewPromptService() *PromptServiceImpl {
	return &PromptServiceImpl{}
}
