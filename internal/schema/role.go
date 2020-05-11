// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type RolePure struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Accession   []string `json:"accession"`
	BuildIn     bool     `json:"build_in"`
	Note        *string  `json:"note"`
}

type Role struct {
	RolePure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
