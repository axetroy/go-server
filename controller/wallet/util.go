package wallet

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/core/errors"
	"gitlab.com/axetroy/server/model"
	"strings"
)

func GetTableName(currency string) string {
	return "wallet_" + strings.ToLower(currency)
}

func EnsureWalletExist(session *xorm.Session, currency string, uid int64) (*model.Wallet, error) {
	var (
		isExistTable  bool
		isExistWallet bool
		tableName     = GetTableName(currency)
		err           error
		w             = model.Wallet{}
	)

	if isExistTable, err = session.IsTableExist(tableName); err != nil {
		return nil, err
	}

	if isExistTable != true {
		err = errors.New(fmt.Sprintf("币种'%v'不存在", currency))
		return nil, err
	}

	query := session.Table(tableName).Where("id = ?", uid)

	if isExistWallet, err = query.Get(&w); err != nil {
		return nil, err
	} else {
		if isExistWallet != true {
			if _, err = session.Insert(&w); err != nil {
				return nil, err
			}

			if isExistWallet, err = query.Get(&w); err != nil {
				return nil, err
			}

			if isExistWallet != true {
				err = errors.New(fmt.Sprintf("用户不存在%v钱包", currency))
				return nil, err
			}
		}
	}

	return &w, nil
}
