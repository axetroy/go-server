package news

import (
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/core/errors"
	"net/http"
)

type CreateNewParams struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Type    model.NewsType `json:"type"`
	Tags    []string       `json:"tags"`
}

type PureNews struct {
	Id      string           `json:"id"`      // 新闻公告类ID
	Author  string           `json:"author"`  // 公告的作者ID
	Tittle  string           `json:"tittle"`  // 公告标题
	Content string           `json:"content"` // 公告内容
	Type    model.NewsType   `json:"type"`    // 公告类型
	Tags    []string         `json:"tags"`    // 公告的标签
	Status  model.NewsStatus `json:"status"`  // 公告状态
}

type Instance struct {
	PureNews
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func Create(uid string, input CreateNewParams) (res response.Response) {
	var (
		err     error
		data    Instance
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

	// TODO: 验证这个UID是不是管理员

	if !model.IsValidNewsType(input.Type) {
		err = exception.NewsInvalidType
	}

	tx = true

	n := model.News{
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
