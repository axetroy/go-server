package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/controller/uploader"
	"github.com/gin-gonic/gin"
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
