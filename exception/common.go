package exception

import "errors"

type exception struct {
	message string
	code    int
}

func New(text string) *exception {
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

var (
	Unknown       = errors.New("未知错误")
	InvalidParams = errors.New("参数不正确")
	NoData        = errors.New("找不到数据")
	InvalidId     = errors.New("ID不正确")
	// auth
	UserNotLogin             = errors.New("请先登陆")
	InvalidAuth              = errors.New("无效的身份认证方式")
	InvalidToken             = errors.New("无效的身份令牌")
	TokenExpired             = errors.New("身份令牌已过期")
	RequirePassword          = errors.New("请输入密码")
	RequirePayPassword       = errors.New("请输入交易密码")
	InvalidPassword          = errors.New("密码错误")
	InvalidAccountOrPassword = errors.New("账号或密码错误")
	InvalidResetCode         = errors.New("密码重置链接已超时")
	InvalidActiveCode        = errors.New("激活链接已超时")
	UserHaveActive           = errors.New("用户已激活")
	PasswordDuplicate        = errors.New("新密码和旧密码不能相同")
	// user
	UserExist    = errors.New("用户已存在")
	UserNotExist = errors.New("用户不存在")
	// upload
	NotSupportType = errors.New("不支持该文件类型")
	OutOfSize      = errors.New("超出文件大小限制")
	// wallet
	NotEnoughBalance = errors.New("钱包余额不足")
)
