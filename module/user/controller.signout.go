// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"github.com/axetroy/go-server/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignOut(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, schema.Response{
		Status:  schema.StatusSuccess,
		Message: "您已登出",
		Data:    true,
	})
}
