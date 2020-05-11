// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAccession() (res schema.Response) {
	var (
		err  error
		data []*accession.Accession
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

	data = accession.List

	return
}

func GetAccessionRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = GetAccession()
}
