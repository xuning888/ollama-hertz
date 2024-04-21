package prompt

import "time"

type PromptPageReq struct {
	PageNum  int    `query:"pageNum,default=0" json:"pageNum"`
	PageSize int    `query:"pageSize,default=10" json:"pageSize"`
	Name     string `query:"name" json:"name"`
}

type PromptPageRes struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Temperature float64   `json:"temperature"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PromptAddOrUpdate struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name,required"`
	Description string  `json:"description,required"`
	Temperature float64 `json:"temperature,default=0.0" vd:"$>=0.0 && $<=2.0; msg:'temperature的值范围[0.0, 2.0]'"`
}
