// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

import "github.com/axetroy/go-server/src/model"

type AdminProfilePure struct {
	Id       string            `json:"id"`       // 用户ID
	Username string            `json:"username"` // 用户名, 用于登陆
	Name     string            `json:"name"`     // 管理员名
	IsSuper  bool              `json:"is_super"` // 是否是超级管理员, 超级管理员全站应该只有一个
	Status   model.AdminStatus `json:"status"`   // 状态
}

type AdminProfileWithToken struct {
	AdminProfile
	Token string `json:"token"`
}

type AdminProfile struct {
	AdminProfilePure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
