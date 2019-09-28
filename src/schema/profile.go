// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

// 用户自己的资料
type ProfilePure struct {
	Id         string   `json:"id"`
	Username   string   `json:"username"`
	Nickname   *string  `json:"nickname"`
	Email      *string  `json:"email"`
	Phone      *string  `json:"phone"`
	Status     int32    `json:"status"`
	Gender     int      `json:"gender"`
	Avatar     string   `json:"avatar"`
	Role       []string `json:"role"`
	Level      int32    `json:"level"`
	InviteCode string   `json:"invite_code"`
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
	PayPassword bool   `json:"pay_password"` // 是否已设置交易密码
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
