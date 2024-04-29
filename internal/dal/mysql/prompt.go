package mysql

import (
	"context"
	"github.com/xuning888/ollama-hertz/internal/model/prompt"
	"github.com/xuning888/ollama-hertz/pkg/logger"
)

type PromptDao struct {
	lg logger.Logger
}

func (pp *PromptDao) PromptPageInfo(ctx context.Context, pageSze int, pageNum int, nameLike string) (
	prompts []*prompt.Prompt, total int64, err error) {

	query := DB.WithContext(ctx).Model(&prompt.Prompt{})
	if nameLike != "" {
		query = query.Where("name like ?", "%"+nameLike+"%")
	}

	if err = query.Count(&total).Error; err != nil {
		pp.lg.Errorf("PromptPageInfo prompt PageInfo total failed with error: %v", err)
		return
	}

	offset := (pageNum - 1) * pageSze
	if err = query.Limit(pageSze).Offset(offset).Find(&prompts).Error; err != nil {
		pp.lg.Errorf("prompt PageInfo limit failed with error: %v", err)
		return
	}
	return
}

func (pp *PromptDao) PromptAdd(ctx context.Context, entity *prompt.Prompt) error {
	return DB.WithContext(ctx).Model(&prompt.Prompt{}).Save(entity).Error
}

func (pp *PromptDao) PromptUpdateById(ctx context.Context, entity *prompt.Prompt) error {
	return DB.WithContext(ctx).Model(&prompt.Prompt{}).Where("id = ?", entity.ID).Updates(entity).Error
}

func NewPromptDao() *PromptDao {
	promptDao := &PromptDao{
		lg: logger.Named("PromptDao"),
	}
	return promptDao
}
