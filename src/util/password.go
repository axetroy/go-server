// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util

const (
	passwordPrefix = "prefix"
	passwordSuffix = "suffix"
)

func GeneratePassword(text string) string {
	password := MD5(passwordPrefix + text + passwordSuffix)
	return password
}
