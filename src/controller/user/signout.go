package user

import (
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignOut(context *gin.Context) {
	context.JSON(http.StatusOK, schema.Response{
		Status:  schema.StatusSuccess,
		Message: "您已登出",
		Data:    true,
	})
}
