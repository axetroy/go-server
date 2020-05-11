// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package schema

type Area struct {
	Province map[string]string `json:"province"`
	City     map[string]string `json:"city"`
	Area     map[string]string `json:"area"`
}
