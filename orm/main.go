package orm

import (
	"fmt"
	"github.com/axetroy/go-server/model"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

var (
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

	fmt.Println("正在同步数据库...")

	db, err := gorm.Open(driverName, DataSourceName)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(
		new(model.Admin),     // 管理员表
		new(model.News),      // 新闻公告
		new(model.User),      // 用户表
		new(model.WalletCny), // 钱包
		new(model.WalletUsd),
		new(model.WalletCoin),
		new(model.InviteHistory),  // 邀请表
		new(model.LoginLog),       // 登陆成功表
		new(model.TransferLogCny), // 钱包转账地址
		new(model.TransferLogUsd),
		new(model.TransferLogCoin),
		new(model.FinanceLogCny), // 流水列表
		new(model.FinanceLogUsd),
		new(model.FinanceLogCoin),
		new(model.Notification), // 系统消息
	)

	DB = db
}
