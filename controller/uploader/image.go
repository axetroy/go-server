package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/response"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

type ImageResponse struct {
	FileResponse
	Thumbnail bool `json:"thumbnail"` // 是否拥有缩略图
}

// 支持的图片后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

func Image(context *gin.Context) {
	var (
		maxUploadSize = Config.Image.MaxSize // 最大上传大小
		distPath      string                 // 最终的输出目录
		err           error
		file          *multipart.FileHeader
		src           multipart.File
		dist          *os.File
		data          *ImageResponse
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if err != nil {
			// TODO: 移除已经下载了的文件
		}

		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusSuccess,
				Message: "",
				Data:    data,
			})
		}
	}()

	// Source
	if file, err = context.FormFile(FIELD); err != nil {
		return
	}

	extname := strings.ToLower(path.Ext(file.Filename))

	if isImage(extname) == false {
		err = exception.NotSupportType
		return
	}

	if maxUploadSize > 0 && file.Size > int64(maxUploadSize) {
		err = exception.OutOfSize
		return
	}

	if src, err = file.Open(); err != nil {
		return
	}
	defer func() {
		if er := src.Close(); er != nil {
			return
		}
	}()

	hash := md5.New()

	if _, err = io.Copy(hash, src); err != nil {
		return
	}

	md5string := hex.EncodeToString(hash.Sum([]byte("")))

	fileName := md5string + extname

	// Destination
	distPath = path.Join(Config.Path, Config.Image.Path, fileName)
	if dist, err = os.Create(distPath); err != nil {
		return
	}
	defer func() {
		if er := dist.Close(); er != nil {
			return
		}
	}()

	// FIXME: open 2 times
	if src, err = file.Open(); err != nil {
		return
	}

	// Copy
	if _, err = io.Copy(dist, src); err != nil {
		return
	}

	var gotThumbnail bool

	// 压缩缩略图
	// 不管成功与否，都会进行下一步的返回
	if _, err := thumbnailify(distPath); err != nil {
		gotThumbnail = false
	} else {
		gotThumbnail = true
	}

	data = &ImageResponse{
		FileResponse: FileResponse{
			Hash:     md5string,
			Filename: fileName,
			Origin:   file.Filename,
			Size:     file.Size,
		},
		Thumbnail: gotThumbnail,
	}
}

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
func thumbnailify(imagePath string) (outputPath string, err error) {
	var (
		file     *os.File
		img      image.Image
		filename = path.Base(imagePath)
	)

	extname := strings.ToLower(path.Ext(imagePath))

	outputPath = path.Join(Config.Path, Config.Image.Thumbnail.Path, filename)

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
		break
	case ".png":
		img, err = png.Decode(file)
		break
	case ".gif":
		img, err = gif.Decode(file)
		break
	default:
		err = exception.NotSupportType
		return
	}

	if img == nil {
		err = errors.New("生成缩略图失败")
		return
	}

	m := resize.Thumbnail(uint(Config.Image.Thumbnail.MaxWidth), uint(Config.Image.Thumbnail.MaxHeight), img, resize.Lanczos3)

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
		break
	case ".png":
		if err = png.Encode(out, m); err != nil {
			return
		}
		break
	case ".gif":
		if err = gif.Encode(out, m, nil); err != nil {
			return
		}
		break
	default:
		err = exception.NotSupportType
		return
	}

	return
}
