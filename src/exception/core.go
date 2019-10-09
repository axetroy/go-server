// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

import (
	"fmt"
	"regexp"
	"strconv"
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

func WrapValidatorError(err error) *Error {
	return New(err.Error(), InvalidParams.Code())
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

	result, err := strconv.Atoi(matchers[1])

	if err != nil {
		return 0
	}

	return result
}
