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

type CreateParams struct {
	Title    string           `json:"title" validate:"required,max=32" comment:"标题"`
	Content  string           `json:"content" validate:"required" comment:"内容"`
	Tags     []string         `json:"tags" validate:"omitempty,max=32" comment:"标签"`
	Status   model.HelpStatus `json:"status" validate:"required" comment:"状态"`
	Type     model.HelpType   `json:"type" validate:"required,oneof=article class" comment:"类型"`
	ParentId *string          `json:"parent_id" validate:"omitempty,max=32" comment:"父级ID"`
}

func Create(c helper.Context, input CreateParams) (res schema.Response) {
	var (
		err  error
		data schema.Help
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

	helpInfo := model.Help{
		Title:    input.Title,
		Content:  input.Content,
		Tags:     input.Tags,
		Status:   input.Status,
		Type:     input.Type,
		ParentId: input.ParentId,
	}

	// checkout parent id is exist or not
	if input.ParentId != nil {
		if err = tx.Where(&model.Help{Id: *input.ParentId}).First(&model.Help{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.HelpParentNotExist
			}
			return
		}
	}

	if err = tx.Create(&helpInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(helpInfo, &data.HelpPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = helpInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = helpInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
