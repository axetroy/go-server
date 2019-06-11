// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/schema"
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
				err = exception.ErrUnknown
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	data = accession.List

	return
}

func GetAccessionRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	res = GetAccession()
}
