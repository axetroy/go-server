// Copyright 2019 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/internal/config"
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
	http.ServeFile(c.Writer, c.Request, originImagePath)
}
