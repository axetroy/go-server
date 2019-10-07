// Copyright 2019 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/config"
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
	http.ServeFile(c.Writer, c.Request, thumbnailImagePath)
}
