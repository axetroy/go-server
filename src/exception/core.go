// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

type exception struct {
	message string
	code    int
}

func New(text string) *exception {
	return &exception{
		message: text,
	}
}

func (e *exception) Error() string {
	return e.message
}

func (e *exception) Code() int {
	return e.code
}
