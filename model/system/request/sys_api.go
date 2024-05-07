package request

type GetApiList struct {
	Page        int    `json:"page" form:"page"`         // 页码
	PagSize     int    `json:"pageSize" form:"pageSize"` // 每页大小
	Path        string `json:"path"`
	Description string `json:"description"`
	ApiGroup    string `json:"apiGroup"`
	Method      string `json:"method"`
}
