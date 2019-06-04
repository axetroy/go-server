// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import "os"

var (
	ModeProduction  = "production"
	ModeDevelopment = "production"
)

type common struct {
	Mode string `json:"mode"`
}

var Common common

func init() {
	if Common.Mode = os.Getenv("GO_MOD"); Common.Mode == "" {
		Common.Mode = ModeDevelopment
	}
}
