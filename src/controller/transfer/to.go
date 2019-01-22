package transfer

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/finance"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

type ToParams struct {
	Currency string  `json:"currency" valid:"required~请选择币种"`                     // 币种
	To       string  `json:"to" valid:"required~请输入转账对象,numeric~请输入正确的接受人ID"`     // 转账给谁
	Amount   string  `json:"amount" valid:"required~请输入转账数量,numeric~请输入纯数字的转账数量"` // 转账数量
	Note     *string `json:"note"`                                                // 转账备注
}

func To(context controller.Context, input ToParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		data         = model.TransferLog{}
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
			fmt.Println(err)
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
		err = exception.InvalidParams
		return
	}

	uid := context.Uid

	tx = service.Db.Begin()

	fromUserInfo := model.User{Id: uid}
	toUserInfo := model.User{Id: input.To}

	if err = tx.Where(&fromUserInfo).Last(&fromUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// TODO: 完善错误信息
			err = exception.UserNotExist
		}
		return
	}

	if err = tx.Where(&toUserInfo).Last(&toUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// TODO: 完善错误信息
			err = exception.UserNotExist
		}
		return
	}

	walletTableName := wallet.GetTableName(input.Currency)      // 对应的钱包表名
	transferTableName := GetTransferTableName(input.Currency)   // 对应的转账记录表名
	financeLogTableName := finance.GetTableName(input.Currency) // 对应的财务日志表名

	fromUserWallet := model.Wallet{
		Id: uid,
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
		if err = tx.CreateTable(model.TransferLogMap[input.Currency]).Error; err != nil {
			return
		}
	}

	transferLog := model.TransferLog{
		Currency: strings.ToUpper(input.Currency),
		From:     uid,
		To:       input.To,
		Status:   model.TransferStatusConfirmed,
		Amount:   amount,
		Note:     input.Note,
	}

	if err = tx.Table(transferTableName).Create(&transferLog).Error; err != nil {
		return
	}

	data = transferLog

	// 如果财务日志表不存在的话, 那么就生成这个表
	if tx.HasTable(financeLogTableName) == false {
		if err = tx.CreateTable(model.FinanceLogMap[input.Currency]).Error; err != nil {
			return
		}
	}

	// 生成我的财务日志
	fromUserFinanceLog := model.FinanceLog{
		OrderId:         transferLog.Id, // 可用余额的变动
		Uid:             uid,
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

func ToRouter(context *gin.Context) {
	var (
		err   error
		input ToParams
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
			fmt.Println(err)
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		return
	}

	res = To(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
