// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.

package model

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/jinzhu/gorm"
	"time"
)

type ConfigField struct {
	Field       string `json:"field"`       // 字段名称
	Description string `json:"description"` // 配置描述
}

var (
	ConfigFieldNamePhone     = ConfigField{Field: "phone", Description: "手机相关的配置"}
	ConfigFieldNameSMTP      = ConfigField{Field: "smtp", Description: "SMTP 邮件服务的配置"}
	ConfigFieldNameWechatApp = ConfigField{Field: "wechat_app", Description: "微信小程序的相关配置"}
	ConfigFields             = []ConfigField{ConfigFieldNamePhone, ConfigFieldNameSMTP, ConfigFieldNameWechatApp}
)

type ConfigFieldPhone struct {
	Provider string `json:"provider"` // 短信服务提供商, 可选 aliyun/tencent
}

type ConfigFieldWechatApp struct {
	AppID  string `json:"app_id" valid:"required~请输入微信小程序的 APP ID"` // 微信小程序的 APP ID
	Secret string `json:"secret" valid:"required~请输入微信小程序的 Secret"` // 微信小程序的密钥
}

type ConfigFieldSMTP struct {
	Server    string `json:"server" valid:"required~请输入 SMTP 地址"`       // SMTP 服务器地址(域名)
	Port      int    `json:"port" valid:"required~请输入 SMTP 端口"`         // SMTP 服务器端口
	Username  string `json:"username" valid:"required~请输入 SMTP 用户名"`    // SMTP 用户名
	Password  string `json:"password" valid:"required~请输入 SMTP 密码"`     // SMTP 密码
	FromName  string `json:"from_name" valid:"required~请输入 SMTP 发送者名字"` // 邮件发送者的名字
	FromEmail string `json:"from_email" valid:"required~请输入 SMTP 邮箱地址"` // 邮件发送者的邮箱地址
}

type Config struct {
	Name      string `gorm:"primary_key;unique;not null;type:varchar(32);index;" json:"name"` // 配置名称
	Fields    string `gorm:"not null;type:text" json:"fields"`                                // 配置对应的字段
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (config *Config) TableName() string {
	return "config"
}

// 校验配置名是否正确
func (config *Config) IsValidConfigName() error {
	// 检验配置名是否正确
	for _, field := range ConfigFields {
		if field.Field == config.Name {
			return nil
		}
	}

	return exception.InvalidParams
}

// 校验配置字段是否正确
func (config *Config) IsValidConfigField() error {
	switch config.Name {
	case ConfigFieldNamePhone.Field:
		c := ConfigFieldPhone{}
		if err := json.Unmarshal([]byte(config.Fields), &c); err != nil {
			return exception.InvalidParams.New(err.Error())
		}
		if err := validator.ValidateStruct(c); err != nil {
			return err
		}
		break
	case ConfigFieldNameSMTP.Field:
		c := ConfigFieldSMTP{}
		if err := json.Unmarshal([]byte(config.Fields), &c); err != nil {
			return exception.InvalidParams.New(err.Error())
		}
		if err := validator.ValidateStruct(c); err != nil {
			return err
		}
		break
	case ConfigFieldNameWechatApp.Field:
		c := ConfigFieldWechatApp{}
		if err := json.Unmarshal([]byte(config.Fields), &c); err != nil {
			return exception.InvalidParams.New(err.Error())
		}
		if err := validator.ValidateStruct(c); err != nil {
			return err
		}
		break
	default:
		return exception.InvalidParams
	}

	return nil
}

func (config *Config) BeforeCreate(scope *gorm.Scope) error {
	if err := config.IsValidConfigName(); err != nil {
		return err
	}

	if err := config.IsValidConfigField(); err != nil {
		return err
	}

	return nil
}
