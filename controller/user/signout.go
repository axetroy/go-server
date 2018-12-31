package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/axetroy/server/response"
	"net/http"
)

func Signout(context *gin.Context) {
	context.JSON(http.StatusOK, response.Response{
		Status:  response.StatusSuccess,
		Message: "您已登出",
		Data:    true,
	})
}