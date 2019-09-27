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
	if Common.Mode = dotenv.Get("GO_MOD"); Common.Mode == "" {
		Common.Mode = ModeProduction
	}
	if Common.MachineId = dotenv.Get("MACHINE_ID"); Common.MachineId == "" {
		Common.MachineId = "0"
	}
	if Common.Signature = dotenv.Get("SIGNATURE_KEY"); Common.Signature == "" {
		Common.Signature = "signature key"
	}
}
