package transfer

import (
	"errors"
	"github.com/axetroy/go-server/controller/finance"
	"github.com/axetroy/go-server/controller/wallet"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
	"strconv"
	"strings"
)

type ToParams struct {
	Currency string  `json:"currency" binding:"required"`    // 币种
	To       string  `json:"to" binding:"required"`          // 转账给谁
	Amount   string  `json:"amount" binding:"required,gt=0"` // 转账数量
	Note     *string `json:"note"`                           // 转账备注
}

func To(context *gin.Context) {
	var (
		err     error
		session *xorm.Session
		tx      bool
		data    = model.TransferLog{}
		input   ToParams
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

	if err = context.ShouldBindJSON(&input); err != nil {
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	myInfo := model.User{Id: uid}
	toInfo := model.User{Id: input.To}

	if isExist, er := session.Get(&myInfo); er != nil {
		return
	} else if isExist != true {
		err = exception.UserNotExist
		return
	}

	if isExist, er := session.Get(&toInfo); er != nil {
		return
	} else if isExist != true {
		err = errors.New("收款人不存在")
		return
	}

	walletTableName := wallet.GetTableName(input.Currency)
	transferTableName := GetTransferTableName(input.Currency)
	financeLogTableName := finance.GetTableName(input.Currency)

	myWallet := model.Wallet{}
	toWallet := model.Wallet{}

	if isExist, er := session.Table(walletTableName).Where("id = ?", uid).Get(&myWallet); er != nil {
		return
	} else if isExist != true {
		err = errors.New("钱包不存在")
		return
	}

	if isExist, er := session.Table(walletTableName).Where("id = ?", input.To).Get(&toWallet); er != nil {
		return
	} else if isExist != true {
		err = errors.New("收款人钱包不存在")
		return
	}

	var amount float64

	if amount, err = strconv.ParseFloat(input.Amount, 64); err != nil {
		return
	}

	if myWallet.Balance < amount {
		err = exception.NotEnoughBalance
		return
	}

	// 变动前的余额/冻结
	myBeforeBalance := myWallet.Balance
	myBeforeFrozen := myWallet.Frozen
	toBeforeBalance := toWallet.Balance
	toBeforeFrozen := toWallet.Frozen

	myWallet.Balance = myWallet.Balance - amount // - 自己的钱包
	toWallet.Balance = toWallet.Balance + amount // + 对方的钱包

	// 变动后的余额/冻结
	myAfterBalance := myWallet.Balance
	myAfterFrozen := myWallet.Frozen
	toAfterBalance := toWallet.Balance
	toAfterFrozen := toWallet.Frozen

	// 余额不能为负数
	if myWallet.Balance < 0 {
		err = exception.NotEnoughBalance
		return
	}

	// 扣除我方的钱
	if _, err = session.Table(walletTableName).Where("id = ?", uid).Cols("balance").ForUpdate().Update(&myWallet); err != nil {
		return
	}

	// 给对方加钱
	if _, err = session.Table(walletTableName).Where("id = ?", input.To).Cols("balance").ForUpdate().Update(&toWallet); err != nil {
		return
	}

	// 如果转账记录的表不存在的话，那么就生成这个表
	if isExist, er := session.IsTableExist(transferTableName); er != nil {
		return
	} else if isExist != true {
		if err = session.CreateTable(model.TransferLogMap[input.Currency]); err != nil {
			return
		}
	}

	transferLog := model.TransferLog{
		Id:       id.Generate(),
		Currency: strings.ToUpper(input.Currency),
		From:     uid,
		To:       input.To,
		Status:   model.TransferStatusConfirmed,
		Amount:   amount,
		Note:     input.Note,
	}

	if _, err = session.Table(transferTableName).Insert(&transferLog); err != nil {
		return
	}

	if _, err = session.Table(transferTableName).Where("id = ?", transferLog.Id).Get(&data); err != nil {
		return
	}

	// ensure finance log table exist
	if isExist, er := session.IsTableExist(financeLogTableName); er != nil {
		return
	} else if isExist != true {
		if err = session.CreateTable(model.FinanceLogMap[input.Currency]); err != nil {
			return
		}
	}

	// 生成我的财务日志
	myFinanceLog := model.FinanceLog{
		Id:      id.Generate(),
		OrderId: transferLog.Id,
		Uid:     uid,
		// 可用余额的变动
		BeforeBalance:   myBeforeBalance,
		BalanceMutation: -amount,
		AfterBalance:    myAfterBalance,
		// 冻结余额的变动
		BeforeFrozen:   myBeforeFrozen,
		FrozenMutation: 0,
		AfterFrozen:    myAfterFrozen,

		Type: model.FinanceTypeTransferOut,
	}

	if _, err = session.Table(financeLogTableName).Insert(&myFinanceLog); err != nil {
		return
	}

	// 生成对方的财务日志
	toFinanceLog := model.FinanceLog{
		Id:      id.Generate(),
		OrderId: transferLog.Id,
		Uid:     input.To,
		// 可用余额的变动
		BeforeBalance:   toBeforeBalance,
		BalanceMutation: amount,
		AfterBalance:    toAfterBalance,
		// 冻结余额的变动
		BeforeFrozen:   toBeforeFrozen,
		FrozenMutation: 0,
		AfterFrozen:    toAfterFrozen,

		Type: model.FinanceTypeTransferIn,
	}

	if _, err = session.Table(financeLogTableName).Insert(&toFinanceLog); err != nil {
		return
	}
}
