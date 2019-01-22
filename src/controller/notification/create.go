package notification

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateParams struct {
	Tittle  string  `json:"tittle"`  // 公告标题
	Content string  `json:"content"` // 公告内容
	Note    *string `json:"note"`    // 备注
}

func Create(context controller.Context, input CreateParams) (res schema.Response) {
	//var (
	//	err     error
	//	data    Notification
	//	session *xorm.Session
	//	tx      bool
	//)
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		switch t := r.(type) {
	//		case string:
	//			err = errors.New(t)
	//		case error:
	//			err = t
	//		default:
	//			err = exception.Unknown
	//		}
	//	}
	//
	//	if tx {
	//		if err != nil {
	//			_ = session.Rollback()
	//		} else {
	//			err = session.Commit()
	//		}
	//	}
	//
	//	if session != nil {
	//		session.Close()
	//	}
	//
	//	if err != nil {
	//		res.Data = nil
	//		res.Message = err.Error()
	//	} else {
	//		res.Data = data
	//		res.Status = response.StatusSuccess
	//	}
	//}()
	//
	//session = orm.Db.NewSession()
	//
	//if err = session.Begin(); err != nil {
	//	return
	//}
	//
	//adminInfo := model.Admin{
	//	Id: context.Uid,
	//}
	//
	//if isExist, er := session.Get(&adminInfo); er != nil {
	//	err = er
	//	return
	//} else if !isExist {
	//	err = exception.AdminNotExist
	//	return
	//}
	//
	//// 需要超级管理员才能创建
	//if !adminInfo.IsSuper {
	//	err = exception.AdminNotSuper
	//	return
	//}
	//
	//tx = true
	//
	//n := model.Notification{
	//	Id:      id.Generate(),
	//	Tittle:  input.Tittle,
	//	Content: input.Content,
	//	Status:  model.NotificationStatusActive,
	//}
	//
	//if _, err = session.Insert(&n); err != nil {
	//	return
	//}
	//
	//fmt.Printf("%+v\n", n)
	//
	//if er := mapstructure.Decode(n, &data.Pure); er != nil {
	//	err = er
	//	return
	//}
	//
	//data.CreatedAt = n.CreatedAt.Format(time.RFC3339Nano)
	//data.UpdatedAt = n.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateRouter(context *gin.Context) {
	var (
		input CreateParams
		err   error
		res   = schema.Response{}
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

	res = Create(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
