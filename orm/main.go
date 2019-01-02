package orm

import (
	"fmt"
	"github.com/axetroy/go-server/model"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

var Db *xorm.Engine

func init() {
	var err error
	if Db, err = xorm.NewEngine("postgres", "postgres://postgres:postgres@localhost:65432/gotest?sslmode=disable"); err != nil {
		panic(err)
		return
	}

	fmt.Println("正在同步数据库...")

	// sync table
	err = Db.Sync(
		// 管理员表
		new(model.Admin),
		// 用户表
		new(model.User),
		new(model.LoginLog),
		// 钱包
		new(model.WalletCny),
		new(model.WalletUsd),
		new(model.WalletCoin),
		// 钱包转账地址
		new(model.TransferLogCny),
		new(model.TransferLogUsd),
		new(model.TransferLogCoin),
		// 流水列表
		new(model.FinanceLogCny),
		new(model.FinanceLogUsd),
		new(model.FinanceLogCoin),
		// 邀请表
		new(model.InviteHistory),
		// 新闻公告表
		new(model.News),
	)

	if err != nil {
		fmt.Println("同步数据库错误")
		log.Fatal(err)
		return
	}

	Db.ShowSQL(true)
}
