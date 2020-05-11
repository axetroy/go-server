// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

// 用户自己的资料
type ProfilePure struct {
	Id                      string   `json:"id"`
	Username                string   `json:"username"`
	Nickname                *string  `json:"nickname"`
	Email                   *string  `json:"email"`
	Phone                   *string  `json:"phone"`
	Status                  int32    `json:"status"`
	Gender                  int      `json:"gender"`
	Avatar                  string   `json:"avatar"`
	Role                    []string `json:"role"`
	Level                   int32    `json:"level"`
	InviteCode              string   `json:"invite_code"`
	UsernameRenameRemaining int      `json:"username_rename_remaining"`
}

// 绑定的微信帐号信息
type WechatBindingInfo struct {
	Nickname  string `json:"nickname"`   // 用户昵称
	AvatarUrl string `json:"avatar_url"` // 用户头像
	Gender    int    `json:"gender"`     // 性别
	Country   string `json:"country"`    // 国家
	Province  string `json:"province"`   // 省份
	City      string `json:"city"`       // 城市
	Language  string `json:"language"`   // 语言
}

// 公开的用户资料，任何人都可以查阅
type ProfilePublic struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Nickname *string `json:"nickname"`
	Avatar   string  `json:"avatar"`
}

type ProfileWithToken struct {
	Profile
	Token string `json:"token"`
}

type Profile struct {
	ProfilePure
	PayPassword bool               `json:"pay_password"` // 是否已设置交易密码
	Wechat      *WechatBindingInfo `json:"wechat"`       // 绑定的微信帐号信息，没有是为 null
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}
