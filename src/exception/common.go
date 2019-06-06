// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	Unknown       = New("未知错误")
	InvalidParams = New("参数不正确")
	NoData        = New("找不到数据")
	InvalidId     = New("ID不正确")
	// auth
	UserNotLogin             = New("请先登陆")
	InvalidAuth              = New("无效的身份认证方式")
	InvalidToken             = New("无效的身份令牌")
	TokenExpired             = New("身份令牌已过期")
	RequirePassword          = New("请输入密码")
	RequirePayPassword       = New("请输入交易密码")
	InvalidPassword          = New("密码错误")
	InvalidAccountOrPassword = New("账号或密码错误")
	InvalidActiveCode        = New("激活链接已超时")
	UserHaveActive           = New("用户已激活")
	PasswordDuplicate        = New("新密码和旧密码不能相同")
	InvalidInviteCode        = New("无效的邀请码")
	PayPasswordSet           = New("交易密码已设置")
	PayPasswordNotSet        = New("请先设置交易密码")
	// user
	UserExist    = New("用户已存在")
	UserNotExist = New("用户不存在")
	// 没有权限
	NoPermission = New("没有权限")

	// 查询列表
	EmptyList = New("sql: no rows in result set")
)
