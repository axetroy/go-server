package controller

type Context struct {
	Uid       string `json:"uid"`
	UserAgent string `json:"user_agent"`
	Ip        string `json:"ip"`
}
