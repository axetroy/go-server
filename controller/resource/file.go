package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/controller/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func File(context *gin.Context) {
	filename := context.Param("filename")
	filePath := path.Join(uploader.Config.Path, uploader.Config.File.Path, filename)
	if isExistFile := fs.PathExists(filePath); isExistFile == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, filePath)
}
