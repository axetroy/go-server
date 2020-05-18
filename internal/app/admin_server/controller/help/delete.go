package help

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

func DeleteHelpById(id string) {
	b := model.Help{}
	database.DeleteRowByTable(b.TableName(), "id", id)
}

func Delete(c helper.Context, id string) (res schema.Response) {
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

	tx = database.Db.Begin()

	adminInfo := model.Admin{Id: c.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	helpInfo := model.Help{
		Id: id,
	}

	if err = tx.First(&helpInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
			return
		}
		return
	}

	if err = tx.Delete(model.Help{
		Id: helpInfo.Id,
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(helpInfo, &data.HelpPure); err != nil {
		return
	}

	data.CreatedAt = helpInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = helpInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var DeleteRouter = router.Handler(func(c router.Context) {
	id := c.Param("help_id")

	c.ResponseFunc(nil, func() schema.Response {
		return Delete(helper.NewContext(&c), id)
	})
})
