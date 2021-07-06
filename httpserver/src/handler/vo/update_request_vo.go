package vo

type UpdateRequestVO struct {
	NickName string `json:"nickName" binding:"required,min=1,max=30"`
}
