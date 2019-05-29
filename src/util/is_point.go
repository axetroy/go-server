package util

import "reflect"

func IsPoint(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return vi.Kind() == reflect.Ptr
}
