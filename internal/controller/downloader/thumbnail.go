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

func Thumbnail(c *gin.Context) {
	filename := c.Param("filename")
	Config := config.Upload
	originImagePath := path.Join(Config.Path, Config.Image.Path, filename)
	thumbnailImagePath := path.Join(Config.Path, Config.Image.Thumbnail.Path, filename)
	if fs.PathExists(thumbnailImagePath) == false {
		// if thumbnail image not exist, try to get origin image
		if fs.PathExists(originImagePath) == true {
			http.ServeFile(c.Writer, c.Request, originImagePath)
			return
		}
		// if the path not found
		http.NotFound(c.Writer, c.Request)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	http.ServeFile(c.Writer, c.Request, thumbnailImagePath)
}
