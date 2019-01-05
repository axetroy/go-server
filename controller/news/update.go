package news

import (
	"errors"
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Tittle  *string           `json:"tittle"`
	Content *string           `json:"content"`
	Type    *model.NewsType   `json:"type"`
	Tags    *[]string         `json:"tags"`
	Status  *model.NewsStatus `json:"status"`
}

func Update(context controller.Context, newsId string, input UpdateParams) (res response.Response) {
	var (
		err          error
		data         News
		session      *xorm.Session
		tx           bool
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

		if tx {
			if err != nil || !shouldUpdate {
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

	{
		if isExist, er := session.Exist(&model.Admin{Id: context.Uid}); er != nil {
			err = er
			return
		} else if !isExist {
			err = exception.AdminNotExist
			return
		}
	}

	newsInfo := model.News{}

	query := session.Where("id = ?", newsId)

	if isExist, er := query.Get(&newsInfo); er != nil {
		err = er
	} else {
		if isExist == false {
			err = exception.NewsNotExist
			return
		}
	}

	if input.Tittle != nil {
		shouldUpdate = true
		newsInfo.Tittle = *input.Tittle
		query = query.Cols("tittle")
	}

	if input.Content != nil {
		shouldUpdate = true
		newsInfo.Content = *input.Content
		query = query.Cols("content")
	}

	if input.Type != nil {
		shouldUpdate = true
		newsInfo.Type = *input.Type
		query = query.Cols("type")
	}

	if input.Status != nil {
		shouldUpdate = true
		newsInfo.Status = *input.Status
		query = query.Cols("status")
	}

	if input.Tags != nil {
		shouldUpdate = true
		newsInfo.Tags = *input.Tags
		query = query.Cols("tags")
	}

	if _, err = query.Update(&newsInfo); err != nil {
		return
	}

	if err = mapstructure.Decode(newsInfo, &data.Pure); err != nil {
		return
	}

	data.CreatedAt = newsInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = newsInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(context *gin.Context) {
	var (
		err   error
		res   = response.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	id := context.Param("news_id")

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.Context{
		Uid: context.GetString("uid"),
	}, id, input)
}
