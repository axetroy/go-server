// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	config2 "github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/router"
	"net/http"
	"path"
)

var Thumbnail = router.Handler(func(c router.Context) {
	filename := c.Param("filename")
	Config := config2.Upload
	originImagePath := path.Join(Config.Path, Config.Image.Path, filename)
	thumbnailImagePath := path.Join(Config.Path, Config.Image.Thumbnail.Path, filename)
	if fs.PathExists(thumbnailImagePath) == false {
		// if thumbnail image not exist, try to get origin image
		if fs.PathExists(originImagePath) == true {
			http.ServeFile(c.Writer(), c.Request(), originImagePath)
			return
		}
		// if the path not found
		http.NotFound(c.Writer(), c.Request())
		return
	}
	http.ServeFile(c.Writer(), c.Request(), thumbnailImagePath)
})
