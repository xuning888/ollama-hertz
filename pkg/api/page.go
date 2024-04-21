package api

type PageInfo struct {
	PageSize int         `json:"pageSize"`
	PageNum  int         `json:"pageNum"`
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
}
