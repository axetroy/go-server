// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
)

func Status(c helper.Context) (res schema.Response) {
	var (
		err  error
		data schema.MessageStatus
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

		helper.Response(&res, data, nil, err)
	}()

	var total int64

	filter := map[string]interface{}{}

	filter["uid"] = c.Uid
	filter["read"] = false

	if err = database.Db.Model(model.Message{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	data.Unread = total

	return
}

// GetRouter get Message detail router
var GetStatusRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return Status(helper.NewContext(&c))
	})
})
