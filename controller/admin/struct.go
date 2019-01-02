package admin

import (
	"github.com/axetroy/go-server/model"
)

type Pure struct {
	Id       string            `json:"id"`       // 用户ID
	Username string            `json:"username"` // 用户名
	Password string            `json:"password"` // 登陆密码
	Name     string            `json:"name"`     // 管理员名
	IsSuper  bool              `json:"is_super"` // 是否是超级管理员, 超级管理员全站应该只有一个
	Status   model.AdminStatus `json:"status"`   // 状态
}

type Detail struct {
	Pure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
