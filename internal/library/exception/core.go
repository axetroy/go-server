// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package exception

import "fmt"

func New(text string, code int) Error {
	return Error{
		message: text,
		code:    code,
	}
}

type Error struct {
	message string
	code    int
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Code() int {
	return e.code
}

func (e Error) New(msg string) Error {
	return New(fmt.Sprintf("%s: %s", e.message, msg), e.code)
}
