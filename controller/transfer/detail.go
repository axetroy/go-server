package transfer

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"net/http"
	"time"
)

func GetDetail(context *gin.Context) {
	var (
		err     error
		session *xorm.Session
		tx      bool
		data    = Log{}
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

		if tx {
			if err != nil {
				_ = session.Rollback()
			} else {
				err = session.Commit()
			}
		}

		if session != nil {
			session.Close()
		}

		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusSuccess,
				Message: "",
				Data:    data,
			})
		}
	}()

	uid := context.GetInt64("uid")

	transferId := context.Param("id")

	if id.IsValidStr(transferId) != true {
		err = exception.InvalidId
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist != true {
		err = exception.UserNotExist
		return
	}

	// 联表查询
	sql := GenerateSql(uid, "*")

	if res, er := session.QueryInterface(sql + " LIMIT 1"); er != nil {
		err = er
		return
	} else {
		if len(res) == 0 {
			err = exception.NoData
			return
		}
		var v = res[0]
		if err = mapstructure.Decode(v, &data); err != nil {
			return
		}
		createdAt := v["created_at"].(time.Time)
		updatedAt := v["updated_at"].(time.Time)

		data.CreatedAt = createdAt.Format(time.RFC3339Nano)
		data.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	}

}
