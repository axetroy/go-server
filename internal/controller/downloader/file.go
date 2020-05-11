// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func File(c *gin.Context) {
	filename := c.Param("filename")
	filePath := path.Join(config.Upload.Path, config.Upload.File.Path, filename)
	if isExistFile := fs.PathExists(filePath); isExistFile == false {
		// if the path not found
		http.NotFound(c.Writer, c.Request)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	http.ServeFile(c.Writer, c.Request, filePath)
}
