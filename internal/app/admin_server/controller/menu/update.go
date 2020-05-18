// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateParams struct {
	Name      *string   `json:"name"`      // 菜单名
	Url       *string   `json:"url"`       // 菜单链接的 URL 地址
	Icon      *string   `json:"icon"`      // 菜单的图标
	Accession *[]string `json:"accession"` // 该菜单所需要的权限
	Sort      *int      `json:"sort"`      // 菜单排序, 越大的越靠前
	ParentId  *string   `json:"parent_id"` // 该菜单的父级 ID
}

func Update(c helper.Context, bannerId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Menu
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err == nil {
			if len(data.Accession) == 0 {
				data.Accession = []string{}
			}
			if len(data.Children) == 0 {
				data.Children = []schema.Menu{}
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
		Id: bannerId,
	}

	if err = tx.First(&menuInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
			return
		}
		return
	}

	var m = map[string]interface{}{}

	if input.Name != nil {
		shouldUpdate = true
		m["name"] = *input.Name
	}

	if input.Url != nil {
		shouldUpdate = true
		m["url"] = *input.Url
	}

	if input.Icon != nil {
		shouldUpdate = true
		m["icon"] = *input.Icon
	}

	if input.Accession != nil {
		// 只保留有效的权限
		shouldUpdate = true
		m["accession"] = accession.FilterAdminAccession(*input.Accession)
	}

	if input.Sort != nil {
		shouldUpdate = true
		m["sort"] = *input.Sort
	}

	if input.ParentId != nil {
		shouldUpdate = true
		m["parent_id"] = *input.ParentId
	}

	if shouldUpdate {
		if err = tx.Model(&menuInfo).Updates(m).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoData
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(menuInfo, &data.MenuPure); err != nil {
		return
	}

	data.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	id := c.Param("menu_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), id, input)
	})
})
