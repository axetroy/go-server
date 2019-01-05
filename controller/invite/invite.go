package invite

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/request"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GenerateInviteCode() string {
	b := make([]byte, 4) // 8 位的邀请码
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(b)
	code := hex.EncodeToString(b)
	return code
}

func GetInviteById(m *model.InviteHistory) (res response.Response) {
	var (
		err     error
		data    Invite
		session *xorm.Session
		tx      bool
		isExist bool
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}
	}()

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	if isExist, err = session.Get(m); err != nil {
		return
	} else if isExist == false {
		err = exception.UserNotExist
		return
	}

	if err = mapstructure.Decode(m, &data.InvitePure); err != nil {
		return
	}

	data.CreatedAt = m.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = m.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetMyInviteList(context *gin.Context) {
	var (
		err     error
		uid     int64
		session *xorm.Session
		tx      bool
		data    = make([]model.InviteHistory, 0)
		meta    = &response.Meta{}
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
			context.JSON(http.StatusOK, response.List{
				Response: response.Response{
					Status:  response.StatusFail,
					Message: err.Error(),
					Data:    nil,
				},
				Meta: nil,
			})
		} else {
			context.JSON(http.StatusOK, response.List{
				Response: response.Response{
					Status:  response.StatusSuccess,
					Message: "",
					Data:    data,
				},
				Meta: meta,
			})
		}
	}()

	query := request.Query{}

	if err = context.BindQuery(&query); err != nil {
		return
	}

	query.Normalize()

	if val, isExist := context.Get("uid"); isExist != true {
		return
	} else {
		if uid, err = strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64); err != nil {
			return
		}
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	var total int64

	// TODO: support sort field
	if total, err = session.Table(model.InviteHistory{}).Where("invitor = ?", uid).
		Limit(query.Limit, query.Limit*query.Page).
		FindAndCount(&data); err != nil {
		return
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

}
