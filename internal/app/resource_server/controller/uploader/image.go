// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"golang.org/x/image/draw"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type ImageResponse struct {
	schema.FileResponse
}

// 支持的图片后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

// src   - source image
// rect  - size we want
// scale - scaler
func scaleTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

var Image = router.Handler(func(c router.Context) {
	var (
		maxUploadSize = config.Upload.Image.MaxSize // 最大上传大小
		err           error
		data          = make([]ImageResponse, 0)
		imageDir      = path.Join(config.Upload.Path, config.Upload.Image.Path)
	)

	defer func() {
		c.JSON(err, data, nil)
	}()

	// Get the max post value size passed via iris.WithPostMaxMemory.
	//maxSize := c.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err = c.Request().ParseMultipartForm(maxUploadSize)

	if err != nil {
		return
	}

	form := c.Request().MultipartForm

	files := form.File["file"]

	// 如果找不到图片
	if len(files) == 0 {
		err = exception.InvalidParams
		return
	}

	for _, file := range files {
		var (
			src  multipart.File // 要读取的文件
			dist *os.File       // 最终输出的文件
		)

		// 判断是否是合法的图片
		extname := strings.ToLower(path.Ext(file.Filename))

		{
			if !isImage(extname) {
				err = exception.NotSupportType
				return
			}

			if maxUploadSize > 0 && file.Size > int64(maxUploadSize) {
				err = exception.OutOfSize
				return
			}
		}

		if src, err = file.Open(); err != nil {
			return
		}

		hash := md5.New()

		if _, err = io.Copy(hash, src); err != nil {
			_ = src.Close()
			return
		} else {
			_ = src.Close()
		}

		md5string := hex.EncodeToString(hash.Sum([]byte("")))

		fileName := md5string + extname

		// 输出到最终文件
		distPath := path.Join(imageDir, fileName)

		if dist, err = os.Create(distPath); err != nil {
			return
		}

		if src, err = file.Open(); err != nil {
			return
		}

		// Copy
		if _, err = io.Copy(dist, src); err != nil {
			_ = src.Close()
			_ = dist.Close()
			return
		} else {
			_ = src.Close()
			_ = dist.Close()
		}

		res := ImageResponse{
			FileResponse: schema.FileResponse{
				Hash:         md5string,
				Filename:     fileName,
				Origin:       file.Filename,
				Size:         file.Size,
				RawPath:      "/v1/resource/image/" + fileName,
				DownloadPath: "/v1/download/image/" + fileName,
			},
		}

		data = append(data, res)
	}
})

/**
check a file is a image or not
*/
func isImage(extName string) bool {
	for i := 0; i < len(supportImageExtNames); i++ {
		if supportImageExtNames[i] == extName {
			return true
		}
	}
	return false
}
