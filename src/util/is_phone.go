package util

import "regexp"

var (
	phoneReg = regexp.MustCompile("^1\\d{10}$")
)

func IsPhone(phoneNumber string) bool {
	return phoneReg.MatchString(phoneNumber)
}
