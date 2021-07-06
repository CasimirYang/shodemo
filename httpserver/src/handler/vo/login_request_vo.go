package vo

type LoginRequestVO struct {
	UserName string `json:"userName" binding:"required,min=1,max=30"`
	Password string `json:"password" binding:"required,min=1,max=30"`
}
