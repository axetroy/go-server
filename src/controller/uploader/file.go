package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func File(context *gin.Context) {
	var (
		isSupportFile bool
		maxUploadSize = Config.File.MaxSize   // 最大上传大小
		allowTypes    = Config.File.AllowType // 可上传的文件类型
		distPath      string                  // 最终的输出目录
		err           error
		file          *multipart.FileHeader
		src           multipart.File
		dist          *os.File
		data          *schema.FileResponse
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
			context.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusSuccess,
				Message: "",
				Data:    data,
			})
		}

	}()

	// Source
	if file, err = context.FormFile("file"); err != nil {
		return
	}

	extname := path.Ext(file.Filename)

	if len(allowTypes) != 0 {
		for i := 0; i < len(allowTypes); i++ {
			if allowTypes[i] == extname {
				isSupportFile = true
				break
			}
		}

		if isSupportFile == false {
			err = exception.NotSupportType
			return
		}
	}

	if file.Size > int64(maxUploadSize) {
		err = exception.OutOfSize
		return
	}

	if src, err = file.Open(); err != nil {
		// open the file fail...
	}
	defer func() {
		if er := src.Close(); er != nil {
			err = er
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
	distPath = path.Join(Config.Path, Config.File.Path, fileName)
	if dist, err = os.Create(distPath); err != nil {
		return
	}
	defer func() {
		if er := dist.Close(); er != nil {
			err = er
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

	data = &schema.FileResponse{
		Hash:     md5string,
		Filename: fileName,
		Origin:   file.Filename,
		Size:     file.Size,
	}
}
