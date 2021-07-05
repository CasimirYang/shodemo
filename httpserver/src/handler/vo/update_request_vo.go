package vo

type UpdateRequestVO struct {
	NickName string `json:"nickName" binding:"required,min=5,max=30"`
}
