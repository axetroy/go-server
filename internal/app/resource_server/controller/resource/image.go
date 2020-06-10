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

// 缩放图片 - 按比例缩放
func scaleImage(img image.Image, scale float64) image.Image {
	if scale == 1 || scale == 0 {
		return img
	}
	// new size of image
	width := float64(img.Bounds().Max.X) * scale
	height := float64(img.Bounds().Max.Y) * scale

	return cropImage(img, int(width), int(height))
}

// 裁剪图片 - 指定尺寸裁剪
func cropImage(img image.Image, width int, height int) image.Image {
	dr := image.Rect(0, 0, int(width), int(height))

	m := scaleTo(img, dr, draw.BiLinear)

	return m
}

// 将文件解析为 img
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

// 将图片写到流中
func EncodeImage(img *image.Image, file *os.File, writer io.Writer, option Query) (err error) {
	extname := path.Ext(file.Name())

	var newImage image.Image

	if option.Scale != nil {
		// 缩放图片
		newImage = scaleImage(*img, *option.Scale)
	} else {
		var (
			width  int
			height int
			i      = *img
		)
		// 裁剪图片
		if option.Width != nil && option.Height != nil {
			width = *option.Width
			height = *option.Height
			newImage = cropImage(*img, width, height)
		} else if option.Width != nil && option.Height == nil {
			// 指定了宽度，那么高度自动
			width = *option.Width

			percent := float64(width) / float64(i.Bounds().Max.X)

			height = int(float64(i.Bounds().Max.Y) * percent)

			newImage = cropImage(*img, width, height)
		} else if option.Width == nil && option.Height != nil {
			// 指定了高度，那么宽度自动
			height = *option.Height

			percent := float64(height) / float64(i.Bounds().Max.Y)

			width = int(float64(i.Bounds().Max.X) * percent)

			newImage = cropImage(*img, width, height)
		} else {
			// 既不指定宽度也不指定高度，则用原图
			newImage = *img
		}
	}

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
	Scale  *float64 `json:"scale" url:"scale" validate:"omitempty,gt=0,lt=1" comment:"缩放比例"` // 缩放比例
	Width  *int     `json:"width" url:"width" validate:"omitempty,gt=0" comment:"宽度"`        // 指定图片的宽度
	Height *int     `json:"height" url:"height" validate:"omitempty,gt=0" comment:"高度"`      // 指定图片的高度
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

	// 如果有错误的，那么就返回原始图片作为 fallback
	defer func() {
		if err != nil {
			http.ServeFile(c.Writer(), c.Request(), originImagePath)
		}
	}()

	// 缩放图片
	if query.Scale != nil {
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

		return
	}

	// 按此村裁剪
	if query.Width != nil || query.Height != nil {
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

		return
	}

	http.ServeFile(c.Writer(), c.Request(), originImagePath)
})
