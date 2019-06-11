// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/finance"
	"github.com/axetroy/go-server/module/finance/finance_model"
	"github.com/axetroy/go-server/module/transfer/transfer_model"
	"github.com/axetroy/go-server/module/transfer/transfer_schema"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/util"
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

func To(context schema.Context, input ToParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		data         = transfer_schema.TransferLog{}
		isValidInput bool
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = common_error.ErrUnknown
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
			res.Status = schema.StatusSuccess
			res.Data = data
		}
	}()

	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	fromUserInfo := user_model.User{Id: context.Uid}
	toUserInfo := user_model.User{Id: input.To}

	if err = tx.Where(&fromUserInfo).Last(&fromUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	if err = tx.Where(&toUserInfo).Last(&toUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	walletTableName := wallet.GetTableName(input.Currency)      // 对应的钱包表名
	transferTableName := GetTransferTableName(input.Currency)   // 对应的转账记录表名
	financeLogTableName := finance.GetTableName(input.Currency) // 对应的财务日志表名

	fromUserWallet := wallet_model.Wallet{
		Id: context.Uid,
	}

	toUserWallet := wallet_model.Wallet{
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
		err = wallet.ErrNotEnoughBalance
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
		err = wallet.ErrNotEnoughBalance
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
		if err = tx.CreateTable(transfer_model.TransferLogMap[strings.ToUpper(input.Currency)]).Error; err != nil {
			return
		}
	}

	transferLog := transfer_model.TransferLog{
		Currency: strings.ToUpper(input.Currency),
		From:     context.Uid,
		To:       input.To,
		Status:   transfer_model.TransferStatusConfirmed,
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
		if err = tx.CreateTable(finance_model.FinanceLogMap[input.Currency]).Error; err != nil {
			return
		}
	}

	// 生成我的财务日志
	fromUserFinanceLog := finance_model.FinanceLog{
		OrderId:         transferLog.Id, // 可用余额的变动
		Uid:             context.Uid,
		BeforeBalance:   fromUserBeforeBalance,
		BalanceMutation: -amount,
		AfterBalance:    fromUserAfterBalance,
		BeforeFrozen:    fromUserBeforeFrozen, // 冻结余额的变动
		FrozenMutation:  0,
		AfterFrozen:     fromUserAfterFrozen,
		Type:            finance_model.FinanceTypeTransferOut,
	}

	// 生成对方的财务日志
	toUserFinanceLog := finance_model.FinanceLog{
		OrderId:         transferLog.Id,
		Uid:             input.To,
		BeforeBalance:   toUserBeforeBalance, // 可用余额的变动
		BalanceMutation: amount,
		AfterBalance:    toUserAfterBalance,
		BeforeFrozen:    toUserBeforeFrozen, // 冻结余额的变动
		FrozenMutation:  0,
		AfterFrozen:     toUserAfterFrozen,
		Type:            finance_model.FinanceTypeTransferIn,
	}

	if err = tx.Table(financeLogTableName).Create(&fromUserFinanceLog).Error; err != nil {
		return
	}

	if err = tx.Table(financeLogTableName).Create(&toUserFinanceLog).Error; err != nil {
		return
	}

	return
}

func ToRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		return
	}

	res = To(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
