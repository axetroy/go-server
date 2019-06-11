// Copyright 2019 Axetroy. All rights reserved. MIT license.
package database

import (
	"fmt"
	"github.com/axetroy/go-server/config"
	"github.com/axetroy/go-server/module/address/address_model"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/finance/finance_model"
	"github.com/axetroy/go-server/module/invite/invite_model"
	"github.com/axetroy/go-server/module/log/log_model"
	"github.com/axetroy/go-server/module/menu/menu_model"
	"github.com/axetroy/go-server/module/message/message_model"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/notification/notification_model"
	"github.com/axetroy/go-server/module/report/report_model"
	"github.com/axetroy/go-server/module/role/role_model"
	"github.com/axetroy/go-server/module/transfer/transfer_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/util"
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
			new(admin_model.Admin),                   // 管理员表
			new(news_model.News),                     // 新闻公告
			new(user_model.User),                     // 用户表
			new(role_model.Role),                     // 角色表 - RBAC
			new(wallet_model.WalletCny),              // 钱包 - CNY
			new(wallet_model.WalletUsd),              // 钱包 - USD
			new(wallet_model.WalletCoin),             // 钱包 - COIN
			new(invite_model.InviteHistory),          // 邀请表
			new(log_model.LoginLog),                  // 登陆成功表
			new(transfer_model.TransferLogCny),       // 转账记录 - CNY
			new(transfer_model.TransferLogUsd),       // 转账记录 - USD
			new(transfer_model.TransferLogCoin),      // 转账记录 - COIN
			new(finance_model.FinanceLogCny),         // 流水列表 - CNY
			new(finance_model.FinanceLogUsd),         // 流水列表 - USD
			new(finance_model.FinanceLogCoin),        // 流水列表 - COIN
			new(notification_model.Notification),     // 系统消息
			new(notification_model.NotificationMark), // 系统消息的已读记录
			new(message_model.Message),               // 个人消息
			new(address_model.Address),               // 收货地址
			new(banner_model.Banner),                 // Banner 表
			new(report_model.Report),                 // 反馈表
			new(menu_model.Menu),                     // 后台管理员菜单
		)

		fmt.Println("数据库同步完成.")
	}

	Db = db

	// 确保超级管理员账号存在
	if err := db.First(&admin_model.Admin{Username: "admin", IsSuper: true}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&admin_model.Admin{
				Username:  "admin",
				Name:      "admin",
				Password:  util.GeneratePassword("admin"),
				Accession: []string{},
				Status:    admin_model.AdminStatusInit,
				IsSuper:   true,
			}).Error

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	defaultRole := role_model.Role{Name: role_model.DefaultUser.Name}

	// 确保有默认的角色
	if err := db.First(&defaultRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&role_model.Role{
				Name:        role_model.DefaultUser.Name,
				Description: role_model.DefaultUser.Description,
				Accession:   role_model.DefaultUser.AccessionArray(),
				BuildIn:     true,
			}).Error
		} else {
			panic(err)
		}
	} else {
		// 如果角色已存在，则同步角色的权限
		if err := db.Model(&defaultRole).Update(&role_model.Role{
			Accession: role_model.DefaultUser.AccessionArray(),
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
