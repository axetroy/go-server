// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type MenuPure struct {
	Id        string   `json:"id"`
	ParentId  string   `json:"parent_id"` // 该菜单的父级 ID
	Name      string   `json:"name"`      // 菜单名
	Url       string   `json:"url"`       // 菜单链接的 URL 地址
	Icon      string   `json:"icon"`      // 菜单的图标
	Accession []string `json:"accession"` // 该菜单所需要的权限
	Sort      int      `json:"sort"`      // 菜单排序, 越大的越靠前
}

type Menu struct {
	MenuPure
	Children  []Menu `json:"children"` // 子菜单
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
