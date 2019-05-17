package transfer

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetDetail(context controller.Context, transferId string) (res schema.Response) {
	var (
		err  error
		tx   *gorm.DB
		data = Log{}
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

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Status = schema.StatusSuccess
			res.Data = data
		}
	}()

	uid := context.Uid

	if util.IsValidIdStr(transferId) != true {
		err = exception.InvalidId
		return
	}

	tx = service.Db.Begin()

	userInfo := model.User{Id: uid}

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 联表查询
	sql := GenerateSql(uid, "*") + " LIMIT 1"

	r := tx.Exec(sql)

	if r.Error != nil {
		return
	}

	// TODO: 解析row

	fmt.Println(r.Row())

	//if res, er := session.QueryInterface(sql + " LIMIT 1"); er != nil {
	//	err = er
	//	return
	//} else {
	//	if len(res) == 0 {
	//		err = exception.NoData
	//		return
	//	}
	//	var v = res[0]
	//	if err = mapstructure.Decode(v, &data); err != nil {
	//		return
	//	}
	//	createdAt := v["created_at"].(time.Time)
	//	updatedAt := v["updated_at"].(time.Time)
	//
	//	data.CreatedAt = createdAt.Format(time.RFC3339Nano)
	//	data.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	//}

	return
}

func GetDetailRouter(context *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetDetail(controller.Context{
		Uid: context.GetString(middleware.ContextUidField),
	}, context.Param("transfer_id"))
}
