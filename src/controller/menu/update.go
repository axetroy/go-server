// Copyright 2019 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/rbac/accession"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
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

func Update(context controller.Context, bannerId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Menu
		tx           *gorm.DB
		shouldUpdate bool
		isValidInput bool
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: context.Uid,
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

func UpdateRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	id := context.Param("menu_id")

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.NewContext(context), id, input)
}
