package vo

type UserResponseVO struct {
	Token    string      `json:"token,omitempty"`
	UserInfo *UserInfoVO `json:"userInfo"`
}

type UserInfoVO struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}
