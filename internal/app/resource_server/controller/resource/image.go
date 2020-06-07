// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	config2 "github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/router"
	"net/http"
	"path"
)

var Image = router.Handler(func(c router.Context) {
	filename := c.Param("filename")
	originImagePath := path.Join(config2.Upload.Path, config2.Upload.Image.Path, filename)
	if !fs.PathExists(originImagePath) {
		// if the path not found
		http.NotFound(c.Writer(), c.Request())
		return
	}
	http.ServeFile(c.Writer(), c.Request(), originImagePath)
})
