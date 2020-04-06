package common

//PushConfig 第三方推送配置
type PushConfig struct {
	URL      string
	TokenURL string
	AppID    string
	Secret   string
}

//UserInfo 用户信息
type UserInfo struct {
	UserID   int
	UserPwd  string
	AppID    string
	Secret   string
	UserType int //用户类型 3 管理用户
}
