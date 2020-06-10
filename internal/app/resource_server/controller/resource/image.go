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

var Image = router.Handler(func(c router.Context) {
	type Query struct {
		Scale  float64 `json:"scale" url:"scale" validate:"omitempty,gt=0,max=1" comment:"缩放比例"` // 缩放比例
		Width  int     `json:"width" url:"with" validate:"omitempty,gt=0" comment:"宽度"`          // 指定图片的宽度
		Height int     `json:"height" url:"height" validate:"omitempty,gt=0" comment:"高度"`       // 指定图片的高度
	}

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

	if err != nil {
		return
	}

	newImage := scaleImage(img, query.Scale)

	c.Header("Content-Type", mime.TypeByExtension(extname))

	switch extname {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(c.Writer(), newImage, nil)
	case ".png":
		err = png.Encode(c.Writer(), newImage)
	case ".gif":
		err = gif.Encode(c.Writer(), newImage, nil)
	default:
		err = exception.NotSupportType
	}

	if err != nil {
		return
	}
})
