package prompt

import "gorm.io/gorm"

type Prompt struct {
	gorm.Model
	// Name 这个提示词的名称是啥
	Name string `json:"name" column:"name"`
	// Description 对这个提示此的描述
	Description string  `json:"description" column:"description"`
	Temperature float64 `json:"temperature" column:"temperature"`
}

func (p *Prompt) TableName() string {
	return "prompt"
}

type Item struct {
	gorm.Model
	PromptId  uint   `json:"promptId" column:"prompt_id"`
	Content   string `json:"content" column:"content"`
	Role      string `json:"role" column:"role"`
	ItemOrder int    `json:"itemOrder" column:"item_order"`
}

func (i *Item) TableName() string {
	return "item"
}
