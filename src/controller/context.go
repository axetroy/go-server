package controller

// 控制器的上下文
type Context struct {
	Uid       string `json:"uid"`        // 操作人
	UserAgent string `json:"user_agent"` // 用户代理
	Ip        string `json:"ip"`         // IP地址
}
