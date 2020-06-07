// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateByAdminParams struct {
	Status *model.ReportStatus `json:"status" validate:"omitempty,number,gte=0" comment:"状态"` // 更改状态
	Locked *bool               `json:"locked" validate:"omitempty" comment:"是否锁定"`            // 是否锁定
}

func UpdateByAdmin(c helper.Context, reportId string, input UpdateByAdminParams) (res schema.Response) {
	var (
		err          error
		data         schema.Report
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

	reportInfo := model.Report{
		Id: reportId,
	}

	if err = tx.First(&reportInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
			return
		}
		return
	}

	updatedModel := map[string]interface{}{}

	if input.Status != nil {
		reportInfo.Status = *input.Status
		updatedModel["status"] = *input.Status
		shouldUpdate = true

		var statusText = ""
		switch *input.Status {
		case 0:
			statusText = "待解决"
		case 1:
			statusText = "已解决"
		default:
			break
		}

		messageInfo := model.Message{
			Uid:   reportInfo.Uid,
			Title: "反馈:" + reportInfo.Title,
			Content: fmt.Sprintf(`感谢你的反馈。
您反馈的 【%s】已被标记为 【%s】
祝您生活愉快
`, reportInfo.Title, statusText),
		}

		// 生成一个用户的个人消息
		if err = tx.Create(&messageInfo).Error; err != nil {
			return
		}

		// 通过 APP 推送给这个用户
		defer func() {
			if err != nil && len(messageInfo.Id) > 0 {
				// 把它加入到队列中
				_ = message_queue.PublishUserMessage(messageInfo.Id)
			}
		}()
	}

	if input.Locked != nil {
		reportInfo.Locked = *input.Locked
		updatedModel["locked"] = *input.Locked
		shouldUpdate = true
	}

	if !shouldUpdate {
		return
	}

	if err = tx.Model(&reportInfo).Where(&model.Report{Id: reportId}).Update(updatedModel).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(reportInfo, &data.ReportPure); err != nil {
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateByAdminRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateByAdminParams
	)

	reportId := c.Param("report_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdateByAdmin(helper.NewContext(&c), reportId, input)
	})
})
