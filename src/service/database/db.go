package database

import (
	"fmt"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

var (
	Db *gorm.DB
)

func init() {
	var (
		err        error
		driverName = os.Getenv("DB_DRIVER")
		dbName     = os.Getenv("DB_NAME")
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbPort     = os.Getenv("DB_PORT")
		sync       = os.Getenv("DB_SYNC")
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

	if len(sync) == 0 {
		sync = "on"
	}

	DataSourceName := fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", driverName, dbUsername, dbPassword, dbPort, dbName)

	fmt.Println("正在连接数据库...")

	db, err := gorm.Open(driverName, DataSourceName)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	if sync == "on" {
		fmt.Println("正在同步数据库...")

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
			new(model.Notification),     // 系统消息
			new(model.NotificationMark), // 系统消息的已读记录
			new(model.Message),          // 个人消息
			new(model.Address),          // 收货地址
			new(model.Banner),           // Banner 表
		)

		fmt.Println("数据库同步完成.")
	}

	Db = db

	// 确保超级管理员账号存在
	if err := db.First(&model.Admin{Username: "admin", IsSuper: true}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Admin{
				Username: "admin",
				Name:     "admin",
				Password: util.GeneratePassword("admin"),
				Status:   model.AdminStatusInit,
				IsSuper:  true,
			}).Error

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func DeleteRowByTable(tableName string, field string, value interface{}) {
	var (
		err error
		tx  *gorm.DB
	)

	defer func() {
		if tx != nil {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}
	}()

	tx = Db.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", tableName, field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}
}
