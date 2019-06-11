// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Status *admin_model.AdminStatus `json:"status"` // 管理员状态
	Name   *string                  `json:"name"`   // 管理员名字
}

func Update(context schema.Context, adminId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         admin_schema.AdminProfile
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
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

	myInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = ErrAdminNotExist
		}
		return
	}

	if !myInfo.IsSuper {
		err = ErrAdminNotSuper
		return
	}

	adminInfo := admin_model.Admin{
		Id: adminId,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrAdminNotExist
			return
		}
		return
	}

	updateModel := admin_model.Admin{}

	if input.Status != nil {
		shouldUpdate = true
		updateModel.Status = *input.Status
		adminInfo.Status = *input.Status
	}

	if input.Name != nil {
		shouldUpdate = true
		updateModel.Name = *input.Name
		adminInfo.Name = *input.Name
	}

	if shouldUpdate {
		if err = tx.Model(&adminInfo).UpdateColumns(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = ErrAdminNotExist
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("admin_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = Update(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id, input)
}
