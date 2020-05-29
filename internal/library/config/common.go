// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type common struct {
	MachineId string `json:"machine_id"` // 机器 ID
	Mode      string `json:"mode"`       // 运行模式, 开发模式还是生产模式
	Exiting   bool   `json:"exiting"`    // 进程是否出于正在退出的状态，用户优雅的退出进程
}

var Common *common

func init() {
	Common = &common{}
	Common.Mode = dotenv.GetByDefault("GO_MOD", "production")
	Common.MachineId = dotenv.GetByDefault("MACHINE_ID", "0")
}
