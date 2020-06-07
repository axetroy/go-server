package validator

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"regexp"
)

var (
	usernameReg = regexp.MustCompile(`^[\w\-]+$`)
	validate    = validator.New()
	trans       ut.Translator
)

func init() {
	z := zh.New()
	uni := ut.New(z, z)
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("zh")
	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatalln(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
}

func ValidateStruct(any interface{}) error {
	if err := validate.Struct(any); err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		errorsMap := errs.Translate(trans)

		msg := ""

		for _, e := range errorsMap {
			if msg == "" {
				msg = e
			} else {
				msg = msg + ";" + e
			}
		}

		return exception.InvalidParams.New(msg)
	}
	return nil
}

func IsEmail(email string) bool {
	err := validate.Var(email, "required,email")

	return err == nil
}

func IsPhone(phone string) bool {
	return util.IsPhone(phone)
}

func IsValidUsername(username string) bool {
	return usernameReg.MatchString(username)
}

func ValidateUsername(username string) error {
	if !IsValidUsername(username) {
		return exception.InvalidFormat
	}
	return nil
}
