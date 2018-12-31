package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/gin-gonic/gin"
	"gitlab.com/axetroy/server/controller/uploader"
	"net/http"
	"path"
)

func Image(context *gin.Context) {
	filename := context.Param("filename")
	originImagePath := path.Join(uploader.Config.Path, uploader.Config.Image.Path, filename)
	if fs.PathExists(originImagePath) == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, originImagePath)
}