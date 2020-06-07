// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	config2 "github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/router"
	"net/http"
	"path"
)

var File = router.Handler(func(c router.Context) {
	filename := c.Param("filename")
	filePath := path.Join(config2.Upload.Path, config2.Upload.File.Path, filename)
	if !fs.PathExists(filePath) {
		// if the path not found
		http.NotFound(c.Writer(), c.Request())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	http.ServeFile(c.Writer(), c.Request(), filePath)
})
