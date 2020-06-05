// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package help

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

type UpdateParams struct {
	Title    *string           `json:"title" validate:"omitempty,max=32" comment:"标题"`
	Content  *string           `json:"content" validate:"omitempty" comment:"内容"`
	Tags     *[]string         `json:"tags" validate:"omitempty,max=32" comment:"标签"`
	Status   *model.HelpStatus `json:"status" validate:"omitempty" comment:"状态"`
	Type     *model.HelpType   `json:"type" validate:"omitempty,oneof=article class" comment:"类型"`
	ParentId *string           `json:"parent_id" validate:"omitempty,max=32" comment:"父级ID"`
}

func Update(c helper.Context, helpId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Help
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

	helpInfo := model.Help{
		Id: helpId,
	}

	if err = tx.First(&helpInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
			return
		}
		return
	}

	updateModel := model.Help{}

	if input.Title != nil {
		shouldUpdate = true
		updateModel.Title = *input.Title
	}

	if input.Content != nil {
		shouldUpdate = true
		updateModel.Content = *input.Content
	}

	if input.Tags != nil {
		shouldUpdate = true
		updateModel.Tags = *input.Tags
	}

	if input.Status != nil {
		shouldUpdate = true
		updateModel.Status = *input.Status
	}

	if input.Type != nil {
		shouldUpdate = true
		updateModel.Type = *input.Type
	}

	if input.ParentId != nil {
		shouldUpdate = true
		updateModel.ParentId = input.ParentId
		// check parent id exist or not
		if er := tx.Model(&helpInfo).Where(model.Help{Id: *input.ParentId}).First(&model.Help{}).Error; er != nil {
			err = er
			return
		}
	}

	if shouldUpdate {
		if err = tx.Model(&helpInfo).Updates(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoData
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(helpInfo, &data.HelpPure); err != nil {
		return
	}

	data.CreatedAt = helpInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = helpInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	id := c.Param("help_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), id, input)
	})
})
