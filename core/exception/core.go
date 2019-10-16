// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Error struct {
	message string
	code    int
}

var (
	errReg = regexp.MustCompile("\\s\\[(\\d+)\\]$")
)

func New(text string, code int) *Error {
	return &Error{
		message: text,
		code:    code,
	}
}

// 包装无效参数错误
func WrapValidatorError(err error) *Error {
	return InheritError(err, InvalidParams)
}

// 继承错误
func InheritError(source error, target *Error) *Error {
	return New(source.Error(), target.Code())
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s [%d]", e.message, e.code)
}

func (e *Error) Code() int {
	return e.code
}

func GetCodeFromError(err error) int {
	msg := err.Error()

	matchers := errReg.FindStringSubmatch(msg)

	if len(matchers) <= 1 {
		return 0
	}

	result, err := strconv.Atoi(strings.TrimSpace(matchers[1]))

	if err != nil {
		return 0
	}

	return result
}
