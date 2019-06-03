package database

import (
	"fmt"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	Db     *gorm.DB
	Config = config.Database
)

func init() {
	DataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Config.Driver, Config.Username, Config.Password, Config.Host, Config.Port, Config.DatabaseName)

	fmt.Println("正在连接数据库...")

	db, err := gorm.Open(Config.Driver, DataSourceName)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	if Config.Sync == "on" {
		fmt.Println("正在同步数据库...")

		// Migrate the schema
		db.AutoMigrate(
			new(model.Admin),            // 管理员表
			new(model.News),             // 新闻公告
			new(model.User),             // 用户表
			new(model.Role),             // 角色表 - RBAC
			new(model.WalletCny),        // 钱包 - CNY
			new(model.WalletUsd),        // 钱包 - USD
			new(model.WalletCoin),       // 钱包 - COIN
			new(model.InviteHistory),    // 邀请表
			new(model.LoginLog),         // 登陆成功表
			new(model.TransferLogCny),   // 转账记录 - CNY
			new(model.TransferLogUsd),   // 转账记录 - USD
			new(model.TransferLogCoin),  // 转账记录 - COIN
			new(model.FinanceLogCny),    // 流水列表 - CNY
			new(model.FinanceLogUsd),    // 流水列表 - USD
			new(model.FinanceLogCoin),   // 流水列表 - COIN
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

	defaultRole := model.Role{Name: model.DefaultUser.Name}

	// 确保有默认的角色
	if err := db.First(&defaultRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Role{
				Name:        model.DefaultUser.Name,
				Description: model.DefaultUser.Description,
				Accession:   model.DefaultUser.AccessionArray(),
				BuildIn:     true,
			}).Error
		} else {
			panic(err)
		}
	} else {
		// 如果角色已存在，则同步角色的权限
		if err := db.Model(&defaultRole).Update(&model.Role{
			Accession: model.DefaultUser.AccessionArray(),
		}).Error; err != nil {
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
