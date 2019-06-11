// Copyright 2019 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/module/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func File(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filePath := path.Join(uploader.Config.Path, uploader.Config.File.Path, filename)
	if isExistFile := fs.PathExists(filePath); isExistFile == false {
		// if the path not found
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	http.ServeFile(ctx.Writer, ctx.Request, filePath)
}
