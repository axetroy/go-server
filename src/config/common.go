// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
)

var (
	ModeProduction = "production"
)

type common struct {
	MachineId string `json:"machine_id"`
	Mode      string `json:"mode"`
	Signature string `json:"signature"`
}

var Common common

func init() {
	Common.Mode = dotenv.GetByDefault("GO_MOD", ModeProduction)
	Common.MachineId = dotenv.GetByDefault("MACHINE_ID", "0")
	Common.Signature = dotenv.GetByDefault("SIGNATURE_KEY", "signature key")
}
