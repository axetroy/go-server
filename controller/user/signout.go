package user

import (
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignOut(context *gin.Context) {
	context.JSON(http.StatusOK, response.Response{
		Status:  response.StatusSuccess,
		Message: "您已登出",
		Data:    true,
	})
}
