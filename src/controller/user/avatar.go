package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/controller/uploader"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type UploadAvatarParams struct {
	Immediately string `form:"immediately"` // 是否立即生效
}

// 支持的头像文件后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png"}

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

func UploadAvatar(uid string, input UploadAvatarParams, file *multipart.FileHeader) (res schema.Response) {
	var (
		err      error
		data     *schema.FileResponse
		tx       *gorm.DB
		src      multipart.File
		distPath string // 最终的输出的文件路径
		dist     *os.File
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	userInfo := model.User{
		Id: uid,
	}

	if src, err = file.Open(); err != nil {
		// open the file fail...
		return
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

	extname := path.Ext(file.Filename)

	if isImage(extname) == false {
		err = exception.NotSupportType
		return
	}

	md5string := hex.EncodeToString(hash.Sum([]byte("")))

	fileName := md5string + extname

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	updateMap := map[string]interface{}{}

	if input.Immediately != "" {
		updateMap["avatar"] = fileName
	}

	if err = tx.Model(&userInfo).Updates(updateMap).Error; err != nil {
		return
	}

	// Destination
	distPath = path.Join(uploader.Config.Path, uploader.Config.Image.Avatar.Path, fileName)

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

	return
}

func UploadAvatarRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UploadAvatarParams
		file  *multipart.FileHeader
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	// Source
	if file, err = context.FormFile("file"); err != nil {
		err = exception.RequireFile
		return
	}

	res = UploadAvatar(context.GetString(middleware.ContextUidField), input, file)
}

func GetAvatarRouter(context *gin.Context) {
	filename := context.Param("filename")
	originImagePath := path.Join(uploader.Config.Path, uploader.Config.Image.Avatar.Path, filename)
	if fs.PathExists(originImagePath) == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, originImagePath)
}
