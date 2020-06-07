// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package uploader

import (
	"crypto/md5"
	"encoding/hex"
	config2 "github.com/axetroy/go-server/internal/app/resource_server/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"io"
	"mime/multipart"
	"os"
	"path"
)

var File = router.Handler(func(c router.Context) {
	var (
		isSupportFile bool
		maxUploadSize = config2.Upload.File.MaxSize   // 最大上传大小
		allowTypes    = config2.Upload.File.AllowType // 可上传的文件类型
		err           error
		data          = make([]schema.FileResponse, 0)
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

				if !isSupportFile {
					err = exception.NotSupportType
					return
				}
			}

			if file.Size > int64(maxUploadSize) {
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
		distPath := path.Join(config2.Upload.Path, config2.Upload.File.Path, fileName)

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

})
