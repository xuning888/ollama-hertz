package service

import (
	"context"
	modelPrompt "github.com/xuning888/yoyoyo/internal/model/prompt"
	"github.com/xuning888/yoyoyo/internal/repo"
	"github.com/xuning888/yoyoyo/internal/schema/prompt"
	"github.com/xuning888/yoyoyo/pkg/logger"
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

type PromptServiceImpl struct {
	promptDao *repo.PromptDao
	lg        logger.Logger
}

func (p *PromptServiceImpl) AddOrUpdate(ctx context.Context, req *prompt.PromptAddOrUpdate) error {
	entity := &modelPrompt.Prompt{}
	entity.ID = req.Id
	entity.Name = req.Name
	entity.Description = req.Description
	entity.Temperature = req.Temperature
	if entity.ID == uint(0) {
		// add
		err := p.promptDao.PromptAdd(ctx, entity)
		if err != nil {
			p.lg.Errorf("AddOrUpdate add prompt failed with error: %v", err)
			return err
		}
		return nil
	}

	err := p.promptDao.PromptUpdateById(ctx, entity)
	if err != nil {
		p.lg.Errorf("AddOrUpdate prompt fialed with error: %v", err)
		return err
	}
	return nil
}

func (p *PromptServiceImpl) PromptPage(ctx context.Context, req *prompt.PromptPageReq) (
	prompts []*prompt.PromptPageRes, total int64, err error) {

	pageSize, pageNum, name := req.PageSize, req.PageNum, req.Name

	dbPrompts, total, err := p.promptDao.PromptPageInfo(ctx, pageSize, pageNum, name)
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

func NewPromptService(promptDao *repo.PromptDao) *PromptServiceImpl {
	return &PromptServiceImpl{
		promptDao: promptDao,
		lg:        logger.Named("PromptServiceImpl"),
	}
}
