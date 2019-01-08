package orm

import (
	"fmt"
	"github.com/axetroy/go-server/env"
	"github.com/axetroy/go-server/model"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	Db *xorm.Engine
	DB *gorm.DB
)

func init() {
	var (
		err        error
		driverName = os.Getenv("DB_DRIVER")
		dbName     = os.Getenv("DB_NAME")
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbPort     = os.Getenv("DB_PORT")
	)

	if len(driverName) == 0 {
		driverName = "postgres"
	}

	if len(dbName) == 0 {
		dbName = "gotest"
	}

	if len(dbUsername) == 0 {
		dbUsername = "postgres"
	}

	if len(dbPassword) == 0 {
		dbPassword = "postgres"
	}

	if len(dbPort) == 0 {
		dbPort = "65432"
	}

	DataSourceName := fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", driverName, dbUsername, dbPassword, dbPort, dbName)

	if Db, err = xorm.NewEngine(driverName, DataSourceName); err != nil {
		panic(err)
	}

	fmt.Println("正在同步数据库...")

	// sync table
	err = Db.Sync(
		new(model.LoginLog),
		// 钱包转账地址
		new(model.TransferLogCny),
		new(model.TransferLogUsd),
		new(model.TransferLogCoin),
		// 流水列表
		new(model.FinanceLogCny),
		new(model.FinanceLogUsd),
		new(model.FinanceLogCoin),
		// 系统消息
		new(model.Notification),
	)

	if err != nil {
		fmt.Println("同步数据库错误")
		log.Fatal(err)
		return
	}

	// 使用 gorm 连接

	db, err := gorm.Open(driverName, DataSourceName)

	if err != nil {
		panic("failed to connect database")
	}

	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(
		// 管理员表
		new(model.Admin),
		// 新闻公告
		new(model.News),
		// 用户表
		new(model.User),
		// 钱包
		new(model.WalletCny),
		new(model.WalletUsd),
		new(model.WalletCoin),
		// 邀请表
		new(model.InviteHistory),
	)
	DB = db

	Db.ShowSQL(env.Test)
}
