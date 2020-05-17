// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package signature

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/schema"
)

func Encryption(c helper.Context, input string) (res schema.Response) {
	var (
		err  error
		data string
		hash string
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	hash, err = util.Signature(input)

	if err != nil {
		return
	}

	data = hash

	return
}

var EncryptionRouter = router.Handler(func(c router.Context) {
	var (
		err   error
		input []byte
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
			c.ResponseFunc(err, func() schema.Response {
				return res
			})
		} else {
			c.ResponseFunc(err, func() schema.Response {
				return Encryption(helper.NewContext(&c), string(input))
			})
		}
	}()

	input, err = c.GetBody()
})
