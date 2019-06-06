// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
)

var (
	ModeProduction  = "production"
	ModeDevelopment = "production"
)

type common struct {
	MachineId string `json:"machine_id"`
	Mode      string `json:"mode"`
}

var Common common

func init() {
	if Common.Mode = dotenv.Get("GO_MOD"); Common.Mode == "" {
		Common.Mode = ModeDevelopment
	}
	if Common.MachineId = dotenv.Get("machine_id"); Common.MachineId == "" {
		Common.MachineId = "0"
	}
}
