// Copyright 2019 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/module/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func Thumbnail(ctx *gin.Context) {
	filename := ctx.Param("filename")
	Config := uploader.Config
	originImagePath := path.Join(Config.Path, Config.Image.Path, filename)
	thumbnailImagePath := path.Join(Config.Path, Config.Image.Thumbnail.Path, filename)
	if fs.PathExists(thumbnailImagePath) == false {
		// if thumbnail image not exist, try to get origin image
		if fs.PathExists(originImagePath) == true {
			http.ServeFile(ctx.Writer, ctx.Request, originImagePath)
			return
		}
		// if the path not found
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	http.ServeFile(ctx.Writer, ctx.Request, thumbnailImagePath)
}
