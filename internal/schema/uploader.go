// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type FileResponse struct {
	Hash         string `json:"hash"`          // 文件 hash
	Filename     string `json:"filename"`      // 存储在服务端的文件名
	Origin       string `json:"origin"`        // 上传文件的原始名
	Size         int64  `json:"size"`          // 文件大小
	RawPath      string `json:"raw_path"`      // 纯文本的文件路径, 需要拼接上域名
	DownloadPath string `json:"download_path"` // 下载的文件路径, 需要拼接上域名
}
