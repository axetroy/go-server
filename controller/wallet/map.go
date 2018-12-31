package wallet

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
)

type Wallets struct {
	Cny  model.Wallet `json:"CNY"`
	Usd  model.Wallet `json:"USD"`
	Coin model.Wallet `json:"COIN"`
}

func GetWallets(context *gin.Context) {
	var (
		err     error
		session *xorm.Session
		tx      bool
		data    = Wallets{}
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

	uid := context.GetString("uid")

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist != true {
		err = exception.UserNotExist
		return
	}

	var (
		cny  *model.Wallet
		usd  *model.Wallet
		coin *model.Wallet
	)

	if cny, err = EnsureWalletExist(session, model.WalletCNY, uid); err != nil {
		return
	}

	if usd, err = EnsureWalletExist(session, model.WalletUSD, uid); err != nil {
		return
	}

	if coin, err = EnsureWalletExist(session, model.WalletCOIN, uid); err != nil {
		return
	}

	data.Cny = *cny
	data.Usd = *usd
	data.Coin = *coin
}
