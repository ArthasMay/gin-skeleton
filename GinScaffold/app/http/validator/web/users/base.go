package users

type Base struct {
	// 必填、对于文本,表示它的长度>=1
	Name string `form:"name" json:"name" binding:"required,min=1"`
}