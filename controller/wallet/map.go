package wallet

import (
	"errors"
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Wallets struct {
	Cny  Wallet `json:"CNY"`
	Usd  Wallet `json:"USD"`
	Coin Wallet `json:"COIN"`
}

func GetWallets(context controller.Context) (res response.Response) {
	var (
		err  error
		data Wallet
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

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}
	}()

	// 获取用户信息
	userInfo := model.User{Id: context.Uid}

	tx = orm.DB.Begin()

	if err = tx.Where(userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// TODO: 如何优雅的union查询

	//sqls := []string{""}

	//for _, v := range model.Wallets {
	//	sql := tx.New().Table("wallet_" + strings.ToLower(v)).Where("id = " + userInfo.Id).QueryExpr()
	//
	//	append(sqls, sql)
	//}
	//
	//orm.DB.Raw("? UNION ?", q1.QueryExpr(), q2.QueryExpr())

	return
}

func GetWalletsRouter(context *gin.Context) {
	var (
		err error
		res = response.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetWallets(controller.Context{
		Uid: context.GetString("uid"),
	})
}
