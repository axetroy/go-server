package wallet

import (
	"github.com/gin-gonic/gin"
)

func GetWallet(context *gin.Context) {
	//var (
	//	err     error
	//	session *xorm.Session
	//	tx      bool
	//	data    = model.Wallet{}
	//)
	//
	//defer func() {
	//
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
	//		context.JSON(http.StatusOK, response.Response{
	//			Status:  response.StatusFail,
	//			Message: err.Error(),
	//			Data:    nil,
	//		})
	//	} else {
	//		context.JSON(http.StatusOK, response.Response{
	//			Status:  response.StatusSuccess,
	//			Message: "",
	//			Data:    data,
	//		})
	//	}
	//}()
	//
	//uid := context.GetString("uid")
	//
	//currency := context.Param("currency")
	//
	//session = orm.Db.NewSession()
	//
	//if err = session.Begin(); err != nil {
	//	return
	//}
	//
	//tx = true
	//
	//defer func() {
	//	if err != nil {
	//		_ = session.Rollback()
	//	} else {
	//		_ = session.Commit()
	//	}
	//}()
	//
	//user := model.User{Id: uid}
	//
	//var isExist bool
	//
	//if isExist, err = session.Get(&user); err != nil {
	//	return
	//}
	//
	//if isExist != true {
	//	err = exception.UserNotExist
	//	return
	//}
	//
	//if w, er := EnsureWalletExist(session, currency, uid); er != nil {
	//	err = er
	//	return
	//} else {
	//	if err = mapstructure.Decode(w, &data); err != nil {
	//		return
	//	}
	//}
}
