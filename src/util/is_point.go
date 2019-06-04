// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util

import "reflect"

func IsPoint(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return vi.Kind() == reflect.Ptr
}
