// Copyright 2019 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/controller/uploader"
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

	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	http.ServeFile(context.Writer, context.Request, filePath)
}
