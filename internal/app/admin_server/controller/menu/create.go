// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type CreateMenuParams struct {
	Name      string    `json:"name" valid:"required~请填写菜单名"` // 菜单名
	Url       *string   `json:"url"`                          // 菜单链接的 URL 地址
	Icon      *string   `json:"icon"`                         // 菜单的图标
	Accession *[]string `json:"accession"`                    // 该菜单所需要的权限
	Sort      *int      `json:"sort" `                        // 菜单排序, 越大的越靠前
	ParentId  *string   `json:"parent_id"`                    // 该菜单的父级 ID
}

type TreeParams struct {
	CreateMenuParams
	Children []TreeParams `json:"children"`
}

func Create(c helper.Context, input CreateMenuParams) (res schema.Response) {
	var (
		err  error
		data schema.Menu
		tx   *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	menuInfo := model.Menu{
		Name: input.Name,
	}

	if input.Url != nil {
		menuInfo.Url = *input.Url
	}

	if input.Icon != nil {
		menuInfo.Icon = *input.Icon
	}

	if input.Sort != nil {
		menuInfo.Sort = *input.Sort
	}

	if input.Accession != nil {
		menuInfo.Accession = *input.Accession
	} else {
		menuInfo.Accession = []string{}
	}

	if input.ParentId != nil {
		menuInfo.ParentId = *input.ParentId

		// 查询是否有这个 parentId
		if err = tx.Where(&model.Menu{Id: *input.ParentId}).Find(&model.Menu{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoData
				return
			}
			return
		}
	}

	if err = tx.Create(&menuInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(menuInfo, &data.MenuPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateMenuParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})

func createChildren(tx *gorm.DB, children []TreeParams, parentId string) ([]*schema.MenuTreeItem, error) {
	var (
		data []*schema.MenuTreeItem
	)

	for _, m := range children {
		menuInfo := model.Menu{
			Name:     m.Name,
			ParentId: parentId,
		}

		if m.Url != nil {
			menuInfo.Url = *m.Url
		}

		if m.Icon != nil {
			menuInfo.Icon = *m.Icon
		}

		if m.Sort != nil {
			menuInfo.Sort = *m.Sort
		}

		if m.Accession != nil {
			menuInfo.Accession = *m.Accession
		} else {
			menuInfo.Accession = []string{}
		}

		if err := tx.Create(&menuInfo).Error; err != nil {
			return nil, err
		}

		info := schema.MenuTreeItem{}

		if err := mapstructure.Decode(menuInfo, &info.MenuPure); err != nil {
			return nil, err
		}

		info.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
		info.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)

		if len(m.Children) > 0 {
			if c, err := createChildren(tx, m.Children, menuInfo.Id); err != nil {
				return nil, err
			} else {
				info.Children = c
			}
		}

		data = append(data, &info)
	}

	return data, nil
}

func CreateFromTree(c helper.Context, input []TreeParams) (res schema.Response) {
	var (
		err  error
		data []schema.MenuTreeItem
		tx   *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	for _, k := range input {
		// 参数校验
		if err = validator.ValidateStruct(k); err != nil {
			return
		}
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	for _, m := range input {
		menuInfo := model.Menu{
			Name: m.Name,
		}

		if m.ParentId != nil {
			menuInfo.ParentId = *m.ParentId
		}

		if m.Url != nil {
			menuInfo.Url = *m.Url
		}

		if m.Icon != nil {
			menuInfo.Icon = *m.Icon
		}

		if m.Sort != nil {
			menuInfo.Sort = *m.Sort
		}

		if m.Accession != nil {
			menuInfo.Accession = *m.Accession
		} else {
			menuInfo.Accession = []string{}
		}

		if err = tx.Create(&menuInfo).Error; err != nil {
			return
		}

		info := schema.MenuTreeItem{}

		if er := mapstructure.Decode(menuInfo, &info.MenuPure); er != nil {
			err = er
			return
		}

		info.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
		info.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)

		if len(m.Children) > 0 {
			if c, err := createChildren(tx, m.Children, menuInfo.Id); err != nil {
				return
			} else {
				var children []*schema.MenuTreeItem

				for _, x := range c {
					children = append(children, x)
				}

				info.Children = children
			}
		}

		data = append(data, info)
	}

	return
}

var CreateFromTreeRouter = router.Handler(func(c router.Context) {
	var (
		input []TreeParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return CreateFromTree(helper.NewContext(&c), input)
	})
})
