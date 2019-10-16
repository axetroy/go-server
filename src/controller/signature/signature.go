// Copyright 2019 Axetroy. All rights reserved. MIT license.
package signature

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Encryption(c controller.Context, input string) (res schema.Response) {
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

		helper.Response(&res, data, err)
	}()

	hash, err = util.Signature(input)

	if err != nil {
		return
	}

	data = hash

	return
}

func EncryptionRouter(c *gin.Context) {
	var (
		err   error
		input []byte
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	input, err = c.GetRawData()

	if err != nil {
		return
	}

	res = Encryption(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, string(input))
}
