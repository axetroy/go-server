// Copyright 2019 Axetroy. All rights reserved. MIT license.
package db

import (
	"fmt"
	"log"

	"github.com/axetroy/go-server/core/config"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/service/dotenv"
	"github.com/axetroy/go-server/core/util"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	connection        *DB
	Config            = config.Database
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func init() {
	DataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Config.Driver, Config.Username, Config.Password, Config.Host, Config.Port, Config.DatabaseName)

	log.Println("正在连接数据库...")

	db, err := gorm.Open(Config.Driver, DataSourceName)

	if err != nil {
		panic(err)
	}

	connection = NewConnection(db)

	connection.EnableLog(config.Common.Mode != "production")

	if Config.Sync == "on" {
		log.Println("正在同步数据库...")

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
			new(model.Report),           // 反馈表
			new(model.Menu),             // 后台管理员菜单
			new(model.Help),             // 帮助中心
			new(model.WechatOpenID),     // 微信 open_id 外键表
			new(model.OAuth),            // oAuth2 表
		)

		log.Println("数据库同步完成.")
	}

	// 确保超级管理员账号存在
	if err := connection.First(&model.Admin{Username: "admin", IsSuper: true}).Error(); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = connection.Create(&model.Admin{
				Username:  "admin",
				Name:      "admin",
				Password:  util.GeneratePassword(dotenv.GetByDefault("ADMIN_DEFAULT_PASSWORD", "admin")),
				Accession: []string{},
				Status:    model.AdminStatusInit,
				IsSuper:   true,
			}).Error()

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	defaultRole := model.Role{Name: model.DefaultUser.Name}

	// 确保有默认的角色
	if err := connection.First(&defaultRole).Error(); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = connection.Create(&model.Role{
				Name:        model.DefaultUser.Name,
				Description: model.DefaultUser.Description,
				Accession:   model.DefaultUser.AccessionArray(),
				BuildIn:     true,
			}).Error()
		} else {
			panic(err)
		}
	} else {
		// 如果角色已存在，则同步角色的权限
		if err := connection.Model(&defaultRole).Update(&model.Role{
			Accession: model.DefaultUser.AccessionArray(),
		}).Error(); err != nil {
			panic(err)
		}
	}

}

type DB struct {
	db *gorm.DB
}

func NewConnection(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) EnableLog(enable bool) *DB {
	d.db = d.db.LogMode(enable)
	return d
}

func (d *DB) Begin() *DB {
	d.db = d.db.Begin()
	return d
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Create(value interface{}) *DB {
	d.db = d.db.Create(value)
	return d
}

func (d *DB) Where(query interface{}, args ...interface{}) *DB {
	d.db = d.db.Where(query, args)
	return d
}

func (d *DB) First(out interface{}, where ...interface{}) *DB {
	d.db = d.db.First(out, where)
	return d
}

func (d *DB) Update(attrs ...interface{}) *DB {
	d.db = d.db.Update(attrs)
	return d
}

func (d *DB) Model(value interface{}) *DB {
	d.db = d.db.Model(value)
	return d
}

func (d *DB) Error() error {
	return d.db.Error
}

func (d *DB) Commit() *DB {
	d.db = d.db.Commit()
	return d
}

func (d *DB) Rollback() *DB {
	d.db = d.db.Rollback()
	return d
}

func DeleteRowByTable(tableName string, field string, value interface{}) {
	var (
		err error
		tx  *DB
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

	tx = connection.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", tableName, field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}
}
