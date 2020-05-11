// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func Image(c *gin.Context) {
	filename := c.Param("filename")
	originImagePath := path.Join(config.Upload.Path, config.Upload.Image.Path, filename)
	if fs.PathExists(originImagePath) == false {
		// if the path not found
		http.NotFound(c.Writer, c.Request)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	http.ServeFile(c.Writer, c.Request, originImagePath)
}
