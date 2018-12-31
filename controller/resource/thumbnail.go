package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/gin-gonic/gin"
	"gitlab.com/axetroy/server/controller/uploader"
	"net/http"
	"path"
)

func Thumbnail(context *gin.Context) {
	filename := context.Param("filename")
	Config := uploader.Config
	originImagePath := path.Join(Config.Path, Config.Image.Path, filename)
	thumbnailImagePath := path.Join(Config.Path, Config.Image.Thumbnail.Path, filename)
	if fs.PathExists(thumbnailImagePath) == false {
		// if thumbnail image not exist, try to get origin image
		if fs.PathExists(originImagePath) == true {
			http.ServeFile(context.Writer, context.Request, originImagePath)
			return
		}
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, thumbnailImagePath)
}
