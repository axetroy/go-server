// Copyright 2019 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/menu/menu_model"
	"github.com/axetroy/go-server/module/menu/menu_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateMenuParams struct {
	Name      string    `json:"name" valid:"required~请填写菜单名"` // 菜单名
	Url       *string   `json:"url"`                          // 菜单链接的 URL 地址
	Icon      *string   `json:"icon"`                         // 菜单的图标
	Accession *[]string `json:"accession"`                    // 该菜单所需要的权限
	Sort      int       `json:"sort" `                        // 菜单排序, 越大的越靠前
	ParentId  *string   `json:"parent_id"`                    // 该菜单的父级 ID
}

func Create(context schema.Context, input CreateMenuParams) (res schema.Response) {
	var (
		err          error
		data         menu_schema.Menu
		tx           *gorm.DB
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	menuInfo := menu_model.Menu{
		Name: input.Name,
	}

	if input.Url != nil {
		menuInfo.Url = *input.Url
	}

	if input.Icon != nil {
		menuInfo.Icon = *input.Icon
	}

	if input.Accession != nil {
		menuInfo.Accession = *input.Accession
	} else {
		menuInfo.Accession = []string{}
	}

	if input.ParentId != nil {
		menuInfo.ParentId = *input.ParentId

		// 查询是否有这个 parentId
		if err = tx.Where(&menu_model.Menu{Id: *input.ParentId}).Find(&menu_model.Menu{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// TODO: 完善错误信息
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

func CreateRouter(ctx *gin.Context) {
	var (
		input CreateMenuParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
