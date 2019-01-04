package news

import (
	"fmt"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/core/errors"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateNewParams struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Type    model.NewsType `json:"type"`
	Tags    []string       `json:"tags"`
}

func Create(uid string, input CreateNewParams) (res response.Response) {
	var (
		err     error
		data    News
		session *xorm.Session
		tx      bool
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}
	}()

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	// TODO: 找一找是否有这个管理员

	// TODO: RBAC权限校验

	if !model.IsValidNewsType(input.Type) {
		err = exception.NewsInvalidType
	}

	tx = true

	n := model.News{
		Id:      id.Generate(),
		Author:  uid,
		Tittle:  input.Title,
		Content: input.Content,
		Type:    input.Type,
		Tags:    input.Tags,
		Status:  model.NewsStatusActive,
	}

	if _, err = session.Insert(&n); err != nil {
		return
	}

	if er := mapstructure.Decode(n, &data.Pure); er != nil {
		err = er
		return
	}

	data.CreatedAt = n.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = n.UpdatedAt.Format(time.RFC3339Nano)

	fmt.Println(n)

	return
}

func CreateRouter(context *gin.Context) {
	var (
		input CreateNewParams
		err   error
		res   = response.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(context.GetString("uid"), input)
}
