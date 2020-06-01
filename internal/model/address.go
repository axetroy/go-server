// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Address struct {
	Id           string  `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 地址ID
	Uid          string  `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 用户ID, 与默认地址联合唯一，用户只能有一个唯一的收货地址
	Name         string  `gorm:"not null;index;type:varchar(32)" json:"name"`                  // 收货人
	Phone        string  `gorm:"not null;index;type:varchar(32)" json:"phone"`                 // 收货人电话
	ProvinceCode string  `gorm:"not null;index;type:varchar(2)" json:"province_code"`          // 省份代码
	CityCode     string  `gorm:"not null;index;type:varchar(4)" json:"city_code"`              // 城市代码
	AreaCode     string  `gorm:"not null;index;type:varchar(6)" json:"area_code"`              // 地区代码
	StreetCode   string  `gorm:"not null;index;type:varchar(9)" json:"street_code"`            // 街道代码
	Address      string  `gorm:"not null;index;type:varchar(32)" json:"address"`               // 详细地址
	IsDefault    bool    `gorm:"not null;index;" json:"is_default"`                            // 是否为默认地址, 跟 UID 联合唯一
	Note         *string `gorm:"null;index;type:varchar(12)" json:"note"`                      // 备注, 通常备注是 家/公司/学校 等
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

func (news *Address) TableName() string {
	return "address"
}

func (news *Address) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
