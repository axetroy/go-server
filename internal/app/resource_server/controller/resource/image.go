// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"golang.org/x/image/draw"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

// src   - source image
// rect  - size we want
// scale - scaler
func scaleTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

// 缩放图片
func scaleImage(img image.Image, scale float64) image.Image {
	if scale == 1 || scale == 0 {
		return img
	}
	// new size of image
	width := float64(img.Bounds().Max.X) * scale
	height := float64(img.Bounds().Max.Y) * scale

	dr := image.Rect(0, 0, int(width), int(height))

	m := scaleTo(img, dr, draw.BiLinear)

	return m
}

// 解码图片
func DecodeImage(file *os.File) (img image.Image, err error) {
	extname := path.Ext(file.Name())
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
	}

	return
}

// 编码图片
func EncodeImage(img *image.Image, file *os.File, writer io.Writer, option Query) (err error) {
	extname := path.Ext(file.Name())

	newImage := scaleImage(*img, option.Scale)

	// decode jpeg into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(writer, newImage, nil)
	case ".png":
		err = png.Encode(writer, newImage)
	case ".gif":
		err = gif.Encode(writer, newImage, nil)
	default:
		err = exception.NotSupportType
	}

	return
}

type Query struct {
	Scale  float64 `json:"scale" url:"scale" validate:"omitempty,gt=0,max=1" comment:"缩放比例"` // 缩放比例
	Width  int     `json:"width" url:"with" validate:"omitempty,gt=0" comment:"宽度"`          // 指定图片的宽度
	Height int     `json:"height" url:"height" validate:"omitempty,gt=0" comment:"高度"`       // 指定图片的高度
}

var Image = router.Handler(func(c router.Context) {
	var (
		img      image.Image
		err      error
		filename = c.Param("filename")
		file     *os.File
		query    Query
	)
	originImagePath := path.Join(config.Upload.Path, config.Upload.Image.Path, filename)
	if !fs.PathExists(originImagePath) {
		// if the path not found
		http.NotFound(c.Writer(), c.Request())
		return
	}

	defer func() {
		if err != nil {
			http.Error(c.Writer(), err.Error(), http.StatusInternalServerError)
		}
	}()

	if err = c.ShouldBindQuery(&query); err != nil {
		return
	}

	if err = validator.ValidateStruct(query); err != nil {
		return
	}

	// 如果不裁剪图片，那么就返回原始图片
	if query.Scale == 0 || query.Scale == 1 {
		http.ServeFile(c.Writer(), c.Request(), originImagePath)
		return
	}

	extname := strings.ToLower(path.Ext(filename))

	// 读取文件
	if file, err = os.Open(originImagePath); err != nil {
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	if img, err = DecodeImage(file); err != nil {
		return
	}

	c.Header("Content-Type", mime.TypeByExtension(extname))

	if err = EncodeImage(&img, file, c.Writer(), query); err != nil {
		return
	}
})
