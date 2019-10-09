// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/finance"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/logger"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/src/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ToParams struct {
	Currency string  `json:"currency" valid:"required~请选择币种"`                   // 币种
	To       string  `json:"to" valid:"required~请输入转账对象,numeric~请输入正确的接受人ID"`   // 转账给谁
	Amount   string  `json:"amount" valid:"required~请输入转账数量,float~请输入纯数字的转账数量"` // 转账数量
	Note     *string `json:"note"`                                              // 转账备注
}

func To(context controller.Context, input ToParams, signature string) (res schema.Response) {
	var (
		err  error
		tx   *gorm.DB
		data = schema.TransferLog{}
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
				logger.Infof("User %s transfer %v", context.Uid, input)
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 交验签名是否正确
	if b, err1 := json.Marshal(input); err != nil {
		err = err1
		return
	} else {
		s, err2 := util.Signature(string(b))

		if err2 != nil {
			err = err2
			return
		}

		// 如果签名不一致
		if s != signature {
			err = exception.InvalidSignature
			return
		}
	}

	tx = database.Db.Begin()

	fromUserInfo := model.User{Id: context.Uid}
	toUserInfo := model.User{Id: input.To}

	if err = tx.Where(&fromUserInfo).Last(&fromUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if err = tx.Where(&toUserInfo).Last(&toUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	walletTableName := wallet.GetTableName(input.Currency)      // 对应的钱包表名
	transferTableName := GetTransferTableName(input.Currency)   // 对应的转账记录表名
	financeLogTableName := finance.GetTableName(input.Currency) // 对应的财务日志表名

	fromUserWallet := model.Wallet{
		Id: context.Uid,
	}

	toUserWallet := model.Wallet{
		Id: input.To,
	}

	if err = tx.Table(walletTableName).Where("id = ?", fromUserWallet.Id).FirstOrInit(&fromUserWallet).Error; err != nil {
		return
	}

	if err = tx.Table(walletTableName).Where("id = ?", toUserWallet.Id).FirstOrInit(&toUserWallet).Error; err != nil {
		return
	}

	var amount float64 // 转账数量

	if amount, err = strconv.ParseFloat(input.Amount, 64); err != nil {
		return
	}

	if fromUserWallet.Balance < amount {
		err = exception.NotEnoughBalance
		return
	}

	// 变动前的余额/冻结
	fromUserBeforeBalance := fromUserWallet.Balance
	fromUserBeforeFrozen := fromUserWallet.Frozen
	toUserBeforeBalance := toUserWallet.Balance
	toUserBeforeFrozen := toUserWallet.Frozen

	fromUserWallet.Balance = fromUserWallet.Balance - amount // - 自己的钱包
	toUserWallet.Balance = toUserWallet.Balance + amount     // + 对方的钱包

	// 变动后的余额/冻结
	fromUserAfterBalance := fromUserWallet.Balance
	fromUserAfterFrozen := fromUserWallet.Frozen
	toUserAfterBalance := toUserWallet.Balance
	toUserAfterFrozen := toUserWallet.Frozen

	// 余额不能为负数
	if fromUserWallet.Balance < 0 {
		err = exception.NotEnoughBalance
		return
	}

	// 扣除我方的钱
	if err = tx.Table(walletTableName).Where("id = ?", fromUserWallet.Id).UpdateColumn("balance", fromUserWallet.Balance).Error; err != nil {
		return
	}

	// 给对方加钱
	if err = tx.Table(walletTableName).Where("id = ?", toUserWallet.Id).UpdateColumn("balance", toUserWallet.Balance).Error; err != nil {
		return
	}

	// 如果转账记录的表不存在的话，那么就生成这个表
	if tx.HasTable(transferTableName) == false {
		if err = tx.CreateTable(model.TransferLogMap[strings.ToUpper(input.Currency)]).Error; err != nil {
			return
		}
	}

	transferLog := model.TransferLog{
		Currency: strings.ToUpper(input.Currency),
		From:     context.Uid,
		To:       input.To,
		Status:   model.TransferStatusConfirmed,
		Amount:   util.FloatToStr(amount), // 保留 8 未小数
		Note:     input.Note,
	}

	if err = tx.Table(transferTableName).Create(&transferLog).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(transferLog, &data.TransferLogPure); err != nil {
		return
	}

	data.CreatedAt = transferLog.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = transferLog.UpdatedAt.Format(time.RFC3339Nano)

	// 如果财务日志表不存在的话, 那么就生成这个表
	if tx.HasTable(financeLogTableName) == false {
		if err = tx.CreateTable(model.FinanceLogMap[input.Currency]).Error; err != nil {
			return
		}
	}

	// 生成我的财务日志
	fromUserFinanceLog := model.FinanceLog{
		OrderId:         transferLog.Id, // 可用余额的变动
		Uid:             context.Uid,
		BeforeBalance:   fromUserBeforeBalance,
		BalanceMutation: -amount,
		AfterBalance:    fromUserAfterBalance,
		BeforeFrozen:    fromUserBeforeFrozen, // 冻结余额的变动
		FrozenMutation:  0,
		AfterFrozen:     fromUserAfterFrozen,
		Type:            model.FinanceTypeTransferOut,
	}

	// 生成对方的财务日志
	toUserFinanceLog := model.FinanceLog{
		OrderId:         transferLog.Id,
		Uid:             input.To,
		BeforeBalance:   toUserBeforeBalance, // 可用余额的变动
		BalanceMutation: amount,
		AfterBalance:    toUserAfterBalance,
		BeforeFrozen:    toUserBeforeFrozen, // 冻结余额的变动
		FrozenMutation:  0,
		AfterFrozen:     toUserAfterFrozen,
		Type:            model.FinanceTypeTransferIn,
	}

	if err = tx.Table(financeLogTableName).Create(&fromUserFinanceLog).Error; err != nil {
		return
	}

	if err = tx.Table(financeLogTableName).Create(&toUserFinanceLog).Error; err != nil {
		return
	}

	return
}

func ToRouter(c *gin.Context) {
	var (
		err   error
		input ToParams
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	// 获取数据签名
	signature := c.GetHeader(middleware.SignatureHeader)

	res = To(controller.NewContext(c), input, signature)
}
