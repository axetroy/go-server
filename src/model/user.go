package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

type UserStatus int32

type Gender int

const (
	// 用户状态
	UserStatusBanned      UserStatus = -100 // 账号被禁用
	UserStatusInactivated UserStatus = -1   // 账号未激活
	UserStatusInit        UserStatus = 0    // 初始化状态

	// 用户性别
	GenderUnknown Gender = 0 // 未知性别
	GenderMale               // 男
	GenderFemmale            // 女
)

type User struct {
	Id            string     `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 用户ID
	Username      string     `gorm:"not null;unique;index" json:"username"`                        // 用户名
	Password      string     `gorm:"not null;type:varchar(36);index" json:"password"`              // 登陆密码
	PayPassword   *string    `gorm:"null;type:varchar(36)" json:"pay_password"`                    // 支付密码
	Nickname      *string    `gorm:"null;type:varchar(36)" json:"nickname"`                        // 昵称
	Phone         *string    `gorm:"null;type:varchar(16);index" json:"phone"`                     // 手机号
	Email         *string    `gorm:"null;type:varchar(36);index" json:"email"`                     // 邮箱
	Status        UserStatus `gorm:"not null" json:"status"`                                       // 状态
	Role          string     `gorm:"not null;type:varchar(36)" json:"role"`                        // 角色
	Avatar        string     `gorm:"not null;type:varchar(36)" json:"avatar"`                      // 头像
	Level         int32      `gorm:"default(1)" json:"level"`                                      // 用户等级
	Gender        Gender     `gorm:"default(0)" json:"gender"`                                     // 性别
	EnableTOTP    bool       `gorm:"not null;" json:"enable_totp"`                                 // 是否启用双重身份认证
	Secret        string     `gorm:"not null;type:varchar(32)" json:"secret"`                      // 用户自己的密钥
	InviteCode    string     `gorm:"not null;unique;type:varchar(8)" json:"invite_code"`           // 用户的邀请码，邀请码唯一
	OauthGoogleId *string    `gorm:"null;unique;type:varchar(255)" json:"oauth_google_id"`         // 用户的GoogleAuth唯一标识符
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

func (news *User) TableName() string {
	return "user"
}

func (news *User) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := util.GenerateId()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}

	// 生成邀请码
	if err := scope.SetColumn("invite_code", util.GenerateInviteCode()); err != nil {
		return err
	}

	// 默认关闭启用谷歌验证码
	if err := scope.SetColumn("enable_totp", false); err != nil {
		return err
	}

	// 生成用户自己的密钥
	if secret, err := util.Generate2FASecret(uid); err != nil {
		return err
	} else {
		if err := scope.SetColumn("secret", secret); err != nil {
			return err
		}
	}

	return nil
}
