package push

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/notify"
)

type CreateNotificationParams struct {
	UserIds []string `json:"user_ids" valid:"required~请输入用户 ID"` // 需要推送的指定用户
	Title   string   `json:"title" valid:"required~请输入标题"`       // 标题
	Content string   `json:"content" valid:"required~请输入内容"`     // 内容
}

func CreateNotification(c helper.Context, input CreateNotificationParams) (res schema.Response) {
	var (
		err error
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 发送通知
	if err := notify.Notify.SendNotifyToCustomUser(input.UserIds, input.Title, input.Content, nil); err != nil {
		err = exception.ThirdParty.New(err.Error())
		return
	}

	return
}

var CreateNotificationRouter = router.Handler(func(c router.Context) {
	var (
		input CreateNotificationParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return CreateNotification(helper.NewContext(&c), input)
	})
})
