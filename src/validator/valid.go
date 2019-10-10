package validator

import (
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/util"
	"regexp"
)

var (
	usernameReg = regexp.MustCompile("^[\\w\\-]+$")
)

func ValidateStruct(any interface{}) error {
	if isValid, err := govalidator.ValidateStruct(any); err != nil {
		return exception.WrapValidatorError(err)
	} else if !isValid {
		return exception.InvalidParams
	}
	return nil
}

func IsEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func IsPhone(phone string) bool {
	return util.IsPhone(phone)
}

func IsValidUsername(username string) bool {
	return usernameReg.MatchString(username)
}

func ValidatePhone(phone string) error {
	if !IsPhone(phone) {
		return exception.InvalidFormat
	} else {
		return nil
	}
}

func ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return exception.InvalidFormat
	} else {
		return nil
	}
}

func ValidateUsername(username string) error {
	if !IsValidUsername(username) {
		return exception.InvalidFormat
	}
	return nil
}
