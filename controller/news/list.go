package news

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/request"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
)

type Query struct {
	request.Query
	Status model.NewsStatus `json:"status" form:"status"`
}

func GetList(input Query) (res response.List) {
	var (
		err     error
		session *xorm.Session
		tx      bool
		data    = make([]model.News, 0)
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
			res.Message = err.Error()
			res.Data = nil
			res.Meta = nil
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
			res.Meta = meta
		}
	}()

	query := input.Query

	query.Normalize()

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	var total int64

	// TODO: support sort field
	if total, err = session.Table(model.News{}).
		Limit(query.Limit, query.Limit*query.Page).
		FindAndCount(&data); err != nil {
		return
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

func GetListRouter(context *gin.Context) {
	var (
		err   error
		res   = response.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetList(input)
}
