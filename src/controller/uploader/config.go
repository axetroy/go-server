// Copyright 2019 Axetroy. All rights reserved. MIT license.
package uploader

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/service/dotenv"
	"path"
)

type FileConfig struct {
	Path      string   `json:"path"`       // 普通文件的存放目录
	MaxSize   int      `json:"max_size"`   // 普通文件上传的限制大小，单位byte, 最大单位1GB
	AllowType []string `json:"allow_type"` // 允许上传的文件后缀名
}

type ImageConfig struct {
	Path      string          `json:"path"`      // 图片存储路径
	MaxSize   int             `json:"max_size"`  // 最大图片上传限制，单位byte
	Thumbnail ThumbnailConfig `json:"thumbnail"` // 缩略图配置
	Avatar    AvatarConfig    `json:"avatar"`    // 用户头像的配置
}

type ThumbnailConfig struct {
	Path      string `json:"path"`       // 缩略图存放路径
	MaxWidth  int    `json:"max_width"`  // 缩略图最大宽度
	MaxHeight int    `json:"max_height"` // 缩略图最大高度
}

type AvatarConfig struct {
	Path string // 头像存储的路径
}

type TConfig struct {
	Path  string      `json:"path"`  //文件上传的根目录
	File  FileConfig  `json:"file"`  // 普通文件上传的配置
	Image ImageConfig `json:"image"` // 普通图片上传的配置
}

var Config = TConfig{
	Path: dotenv.Get("UPLOAD_DIR"),
	File: FileConfig{
		Path:    "file",
		MaxSize: 1024 * 1024 * 10, // max 10MB
	},
	Image: ImageConfig{
		Path:    "image",
		MaxSize: 1024 * 1024 * 10, // max 10MB
		Thumbnail: ThumbnailConfig{
			Path:      "thumbnail",
			MaxWidth:  100,
			MaxHeight: 100,
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
		rootPath = dotenv.Get("UPLOAD_DIR")
	)

	if rootPath == "" {
		Config.Path = path.Join(dotenv.RootDir, "upload")
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
