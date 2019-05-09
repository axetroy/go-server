package uploader

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/util"
	"os"
	"path"
)

type FileConfig struct {
	Path      string   `binding:"required,length(1|20)" json:"path"`        // 普通文件的存放目录
	MaxSize   int      `binding:"required" json:"max_size"`                 // 普通文件上传的限制大小，单位byte, 最大单位1GB
	AllowType []string `binding:"required,length(0|100)" json:"allow_type"` // 允许上传的文件后缀名
}

type ImageConfig struct {
	Path      string          `valid:"required,length(1|20)" json:"path"` // 图片存储路径
	MaxSize   int             `valid:"required" json:"max_size"`          // 最大图片上传限制，单位byte
	Thumbnail ThumbnailConfig // 缩略图配置
	Avatar    AvatarConfig    // 用户头像的配置
}

type ThumbnailConfig struct {
	Path      string `valid:"required,length(1|20)" json:"path"` // 缩略图存放路径
	MaxWidth  int    `valid:"required" json:"max_width"`         // 缩略图最大宽度
	MaxHeight int    `valid:"required" json:"max_height"`        // 缩略图最大高度
}

type AvatarConfig struct {
	Path string // 头像存储的路径
}

type TConfig struct {
	Path  string      `valid:"required,length(1|20)"` //文件上传的根目录
	File  FileConfig  // 普通文件上传的配置
	Image ImageConfig // 普通图片上传的配置
}

var Config = TConfig{
	Path: os.Getenv("UPLOAD_DIR"),
	File: FileConfig{
		Path:    "file",
		MaxSize: 1024 * 1024 * 10, // max 10MB
	},
	Image: ImageConfig{
		Path:    "image",
		MaxSize: 1024 * 1024 * 10, // max 10MB
		Thumbnail: ThumbnailConfig{
			Path:      "thumbnail",
			MaxWidth:  60,
			MaxHeight: 60,
		},
		Avatar: AvatarConfig{
			Path: "avatar",
		},
	},
}

// 确保上传的文件目录存在
func init() {
	var (
		err      error
		cwd      string
		rootPath = os.Getenv("UPLOAD_DIR")
	)

	if len(rootPath) == 0 {
		// 如果是测试环境的话, 写死生成的目录，防止到处创建upload目录
		if util.Test {
			Config.Path = path.Join(util.RootDir, "upload")
		} else if cwd, err = os.Getwd(); err != nil {
			panic(cwd)
		} else {
			Config.Path = path.Join(cwd, "upload")
		}
	} else {
		Config.Path = rootPath
	}

	if err = fs.EnsureDir(path.Join(Config.Path, Config.File.Path)); err != nil {
		return
	}

	if err = fs.EnsureDir(path.Join(Config.Path, Config.Image.Path)); err != nil {
		return
	}

	if err = fs.EnsureDir(path.Join(Config.Path, Config.Image.Thumbnail.Path)); err != nil {
		return
	}

	if err = fs.EnsureDir(path.Join(Config.Path, Config.Image.Avatar.Path)); err != nil {
		return
	}

	return
}
