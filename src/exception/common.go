// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	SystemMaintenance = New("系统维护中", -1)
	Unknown           = New("未知错误", 0)
	InvalidParams     = New("参数不正确", 100000)
	NoData            = New("找不到数据", 100001)
	NoPermission      = New("没有权限", 100002)
	InvalidSignature  = New("数据签名不正确", 100003)
	InvalidFormat     = New("格式不正确", 100004)
	InvalidInviteCode = New("无效的邀请码", 100005)
	SendMsgFail       = New("发送短信失败", 101000)
	SendEmailFail     = New("发送邮件失败", 101001)
	UserNotLogin      = New("请先登陆", 999999)
	InvalidAuth       = New("无效的身份认证方式", 999999)
	InvalidToken      = New("无效的身份令牌", 999999)
	TokenExpired      = New("身份令牌已过期", 999999)
	EmptyList         = New("sql: no rows in result set", 0)

	// 用户类
	UserNotExist             = New("用户不存在", 200000)
	UserExist                = New("用户已存在", 200001)
	UserHaveActive           = New("用户已激活", 200002)
	UserIsInActive           = New("帐号未激活", 200003)
	UserHaveBeenBan          = New("帐号已被禁用", 200004)
	PasswordDuplicate        = New("新密码和旧密码不能相同", 200005)
	InvalidAccountOrPassword = New("账号或密码错误", 200006)
	InvalidResetCode         = New("重置码错误或已失效", 200007)
	RequirePayPasswordSet    = New("需要先设置交易密码", 200008)
	PayPasswordSet           = New("交易密码已设置", 200009)
	InvalidConfirmPassword   = New("两次输入密码不一致", 200010)
	InvalidOldPassword       = New("旧密码错误", 200011)
	InvalidPassword          = New("密码错误", 200012)
	RequirePassword          = New("请输入密码", 200013)
	RequirePayPassword       = New("请输入交易密码", 200014)
	DuplicateBinding         = New("帐号重复绑定", 200015)
	RenameUserNameFail       = New("无法重命名用户名", 200016)

	// 钱包
	NotEnoughBalance = New("钱包余额不足", 0)
	InvalidWallet    = New("无效的钱包", 0)

	// 上传
	RequireFile    = New("请上传文件", 0)
	NotSupportType = New("不支持该文件类型", 0)
	OutOfSize      = New("超出文件大小限制", 0)

	// 地址
	AddressDefaultNotExist     = New("默认地址不存在", 0)
	AddressNotExist            = New("地址记录不存在", 0)
	AddressInvalidProvinceCode = New("无效的省份代码", 0)
	AddressInvalidCityCode     = New("无效的城市代码", 0)
	AddressInvalidAreaCode     = New("无效的地区代码", 0)

	// 管理员
	AdminExist    = New("管理员已存在", 0)
	AdminNotExist = New("管理员不存在", 0)
	AdminNotSuper = New("只有超级管理员才能操作", 0)

	// banner
	BannerInvalidPlatform = New("无效的平台", 0)
	BannerNotExist        = New("不存在横幅", 0)

	// 帮助中心
	HelpParentNotExist = New("父级不存在", 0)

	// 邀请
	InviteNotExist = New("邀请记录不存在", 0)

	// RBAC 角色
	RoleNotExist     = New("角色不存在", 0)
	RoleCannotUpdate = New("无法更新角色", 0)
	RoleHadBeenUsed  = New("角色正在被使用，无法删除", 0)

	// 系统通知
	NotificationNotExist = New("系统通知不存在", 0)

	// 用户消息
	MessageNotExist = New("用户消息不存在", 0)

	// 新闻资讯
	NewsInvalidType = New("错误的文章类型", 0)
	NewsNotExist    = New("文章不存在", 0)
)
