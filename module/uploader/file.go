// Copyright 2019 Axetroy. All rights reserved. MIT license.
package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/schema"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func File(ctx *gin.Context) {
	var (
		isSupportFile bool
		maxUploadSize = Config.File.MaxSize   // 最大上传大小
		allowTypes    = Config.File.AllowType // 可上传的文件类型
		err           error
		data          = make([]schema.FileResponse, 0)
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.ErrUnknown
			}
		}

		if err != nil {
			ctx.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			ctx.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusSuccess,
				Message: "",
				Data:    data,
			})
		}

	}()

	form, er := ctx.MultipartForm()

	if er != nil {
		err = er
		return
	}

	files := form.File["file"]

	// 不管成功与否，都移除已下载到本地的缓存图片
	defer func() {
		_ = form.RemoveAll()
	}()

	for _, file := range files {
		var (
			src  multipart.File // 要读取的文件
			dist *os.File       // 最终输出的文件
		)
		extname := path.Ext(file.Filename)

		// 判断是否是合法的上传文件
		{
			if len(allowTypes) != 0 {
				for i := 0; i < len(allowTypes); i++ {
					if allowTypes[i] == extname {
						isSupportFile = true
						break
					}
				}

				if isSupportFile == false {
					err = exception.ErrNotSupportType
					return
				}
			}

			if file.Size > int64(maxUploadSize) {
				err = exception.ErrOutOfSize
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
		distPath := path.Join(Config.Path, Config.File.Path, fileName)

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

		res := schema.FileResponse{
			Hash:         md5string,
			Filename:     fileName,
			Origin:       file.Filename,
			Size:         file.Size,
			RawPath:      "/v1/resource/file/" + fileName,
			DownloadPath: "/v1/download/file/" + fileName,
		}

		data = append(data, res)

	}

}
