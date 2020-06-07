// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	config2 "github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type ImageResponse struct {
	schema.FileResponse
	Thumbnail     bool   `json:"thumbnail"`      // 是否拥有缩略图
	ThumbnailPath string `json:"thumbnail_path"` // 缩略图的路径
}

// 支持的图片后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

var Image = router.Handler(func(c router.Context) {
	var (
		maxUploadSize = config2.Upload.Image.MaxSize // 最大上传大小
		err           error
		data          = make([]ImageResponse, 0)
		imageDir      = path.Join(config2.Upload.Path, config2.Upload.Image.Path)
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
			Thumbnail: false,
		}

		// 压缩缩略图
		// 不管成功与否，都会进行下一步的返回
		if _, er := GenerateThumbnail(distPath); er == nil {
			res.Thumbnail = true
			res.ThumbnailPath = "/v1/resource/thumbnail/" + fileName
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

/**
Generate thumbnail
*/
func GenerateThumbnail(imagePath string) (outputPath string, err error) {
	var (
		file         *os.File
		img          image.Image
		filename     = path.Base(imagePath)
		maxWidth     = config2.Upload.Image.Thumbnail.MaxWidth
		maxHeight    = config2.Upload.Image.Thumbnail.MaxHeight
		thumbnailDir = path.Join(config2.Upload.Path, config2.Upload.Image.Thumbnail.Path)
	)

	extname := strings.ToLower(path.Ext(imagePath))

	outputPath = path.Join(thumbnailDir, filename)

	// 读取文件
	if file, err = os.Open(imagePath); err != nil {
		return
	}

	defer func() {
		if er := file.Close(); er != nil {
			err = er
			return
		}
	}()

	// decode jpeg into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".gif":
		img, err = gif.Decode(file)
	default:
		err = exception.NotSupportType
		return
	}

	if img == nil {
		err = errors.New("生成缩略图失败")
		return
	}

	m := resize.Thumbnail(uint(maxWidth), uint(maxHeight), img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer func() {
		if er := out.Close(); er != nil {
			return
		}
	}()

	// write new image to file

	// decode jpeg/png/gif into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		if err = jpeg.Encode(out, m, nil); err != nil {
			return
		}
	case ".png":
		if err = png.Encode(out, m); err != nil {
			return
		}
	case ".gif":
		if err = gif.Encode(out, m, nil); err != nil {
			return
		}
	default:
		err = exception.NotSupportType
		return
	}

	return
}
