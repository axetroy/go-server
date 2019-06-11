// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	ErrUnknown       = NewError("未知错误")
	ErrInvalidParams = NewError("参数不正确")
	ErrNoData        = NewError("找不到数据")
	// auth
	ErrUserNotLogin             = NewError("请先登陆")
	ErrRequirePassword          = NewError("请输入密码")
	ErrRequirePayPassword       = NewError("请输入交易密码")
	ErrInvalidPassword          = NewError("密码错误")
	ErrInvalidAccountOrPassword = NewError("账号或密码错误")
	ErrInvalidActiveCode        = NewError("激活链接已超时")
	ErrUserHaveActive           = NewError("用户已激活")
	ErrPasswordDuplicate        = NewError("新密码和旧密码不能相同")
	ErrInvalidInviteCode        = NewError("无效的邀请码")
	ErrPayPasswordSet           = NewError("交易密码已设置")
	ErrPayPasswordNotSet        = NewError("请先设置交易密码")
	// user
	ErrUserNotExist = NewError("用户不存在")
	// 没有权限
	ErrNoPermission = NewError("没有权限")

	// 查询列表
	ErrEmptyList = NewError("sql: no rows in result set")

	ErrRequireFile    = NewError("请上传文件")
	ErrNotSupportType = NewError("不支持该文件类型")
	ErrOutOfSize      = NewError("超出文件大小限制")
)

type exception struct {
	message string
	code    int
}

func NewError(text string) *exception {
	return &exception{
		message: text,
	}
}

func (e *exception) Error() string {
	return e.message
}

func (e *exception) Code() int {
	return e.code
}
