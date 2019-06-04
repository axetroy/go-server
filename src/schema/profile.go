// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type ProfilePure struct {
	Id         string   `json:"id"`
	Username   string   `json:"username"`
	Nickname   *string  `json:"nickname"`
	Email      *string  `json:"email"`
	Phone      *string  `json:"phone"`
	Status     int32    `json:"status"`
	Avatar     string   `json:"avatar"`
	Role       []string `json:"role"`
	Level      int32    `json:"level"`
	InviteCode string   `json:"invite_code"`
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
