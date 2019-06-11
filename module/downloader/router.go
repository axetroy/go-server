// Copyright 2019 Axetroy. All rights reserved. MIT license.

package downloader

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {

	r.GET("/download/file/:filename", File)           // 下载文件
	r.GET("/download/image/:filename", Image)         // 下载图片
	r.GET("/download/thumbnail/:filename", Thumbnail) // 下载缩略图

	return r
}
