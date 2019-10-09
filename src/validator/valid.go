package validator

import (
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/util"
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

func ValidatePhone(phone string) error {
	if !IsPhone(phone) {
		return exception.InvalidParams
	} else {
		return nil
	}
}

func ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return exception.InvalidParams
	} else {
		return nil
	}
}
