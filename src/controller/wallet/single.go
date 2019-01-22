package wallet

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
	"time"
)

func GetWallet(context controller.Context, currencyName string) (res schema.Response) {
	var (
		err  error
		data schema.Wallet
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
			res.Status = schema.StatusSuccess
		}
	}()

	// 获取用户信息
	userInfo := model.User{Id: context.Uid}

	tx = service.Db.Begin()

	if err = tx.Where(userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	walletInfo := model.Wallet{
		Id: userInfo.Id,
	}

	// TODO: 校验currencyName是否合法

	if err = tx.Table("wallet_" + strings.ToLower(currencyName)).Scan(&walletInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.WalletPure); err != nil {
		return
	}

	data.CreatedAt = walletInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = walletInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetWalletRouter(context *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetWallet(controller.Context{
		Uid: context.GetString("uid"),
	}, context.Param("currency"))
}
