// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package database

import (
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/rbac/role"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

var (
	Db     *gorm.DB
	Config = config.Database
)

func Dispose() {
	if Db != nil {
		_ = Db.Close()
	}
}

func init() {
	if dotenv.Test {
		Connect()
	}
}

func Migrate(db *gorm.DB) error {
	DataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Config.Driver, Config.Username, Config.Password, Config.Host, Config.Port, Config.DatabaseName)

	log.Println("正在连接数据库...")

	if db == nil {
		var err error
		db, err = gorm.Open(Config.Driver, DataSourceName)

		if err != nil {
			return err
		}
	}

	db.LogMode(config.Common.Mode != "production")

	// Migrate the schema
	if err := db.AutoMigrate(
		new(model.Config),              // 配置表
		new(model.Admin),               // 管理员表
		new(model.News),                // 新闻公告
		new(model.Role),                // 角色表 - RBAC
		new(model.User),                // 用户表
		new(model.WalletCny),           // 钱包 - CNY
		new(model.WalletUsd),           // 钱包 - USD
		new(model.WalletCoin),          // 钱包 - COIN
		new(model.InviteHistory),       // 邀请表
		new(model.LoginLog),            // 登陆成功表
		new(model.TransferLogCny),      // 转账记录 - CNY
		new(model.TransferLogUsd),      // 转账记录 - USD
		new(model.TransferLogCoin),     // 转账记录 - COIN
		new(model.FinanceLogCny),       // 流水列表 - CNY
		new(model.FinanceLogUsd),       // 流水列表 - USD
		new(model.FinanceLogCoin),      // 流水列表 - COIN
		new(model.Notification),        // 系统消息
		new(model.NotificationMark),    // 系统消息的已读记录
		new(model.Message),             // 个人消息
		new(model.Address),             // 收货地址
		new(model.Banner),              // Banner 表
		new(model.Report),              // 反馈表
		new(model.Menu),                // 后台管理员菜单
		new(model.Help),                // 帮助中心
		new(model.WechatOpenID),        // 微信 open_id 外键表
		new(model.OAuth),               // oAuth2 表
		new(model.CustomerSession),     // 客服会话表
		new(model.CustomerSessionItem), // 客服会话内容表
	).Error; err != nil {
		return err
	}

	log.Println("数据库同步完成.")

	superAdminInfo := model.Admin{Username: "admin", IsSuper: true}

	// 确保超级管理员账号存在
	if err := db.Where(&superAdminInfo).First(&superAdminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Admin{
				Username:  "admin",
				Name:      "admin",
				Password:  util.GeneratePassword(dotenv.GetByDefault("ADMIN_DEFAULT_PASSWORD", "123456")),
				Accession: []string{},
				Status:    model.AdminStatusInit,
				IsSuper:   true,
			}).Error

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 确保内置的角色存在
	buildInRoles := []*role.Role{
		model.DefaultUser,
		model.DefaultWaiter,
	}

	for _, buildInRole := range buildInRoles {
		defaultRole := model.Role{Name: buildInRole.Name}

		// 确保有默认的角色
		if err := db.First(&defaultRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err = db.Create(&model.Role{
					Name:        buildInRole.Name,
					Description: buildInRole.Description,
					Accession:   buildInRole.AccessionArray(),
					BuildIn:     true,
				}).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// 如果角色已存在，则同步角色的权限
			if err := db.Model(&defaultRole).Update(&model.Role{
				Accession: buildInRole.AccessionArray(),
				BuildIn:   true,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func Connect() {
	DataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Config.Driver, Config.Username, Config.Password, Config.Host, Config.Port, Config.DatabaseName)

	log.Println("正在连接数据库...")

	db, err := gorm.Open(Config.Driver, DataSourceName)

	if err != nil {
		log.Fatalln(err)
	}

	db.LogMode(config.Common.Mode != "production")

	Db = db
}

// WANING: 该操作会删除数据，并且不可恢复
// 通常只用于测试中
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

	raw := fmt.Sprintf("DELETE FROM \"%s\" WHERE %s = '%s'", tableName, field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}
}
