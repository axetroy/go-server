// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package schema

import (
	"encoding/json"
	"errors"

	"github.com/axetroy/go-server/internal/library/util"
)

type Meta struct {
	Limit    int     `json:"limit"`              // 当前请求获取多少条数据， 默认 10
	Page     int     `json:"page"`               // 当前第几页，默认 0 开始
	Total    int64   `json:"total"`              // 数据总量
	Num      int     `json:"num"`                // 当前返回的数据流
	Sort     string  `json:"sort"`               // 排序
	Platform *string `json:"platform,omitempty"` // 平台
}

const (
	StatusSuccess = 1
	StatusFail    = 0
)

type Response struct {
	Message string      `json:"message"`        // 附带的消息，接口请求错误时，一般都会有错误信息
	Data    interface{} `json:"data"`           // 接口附带的数据
	Status  int         `json:"status"`         // 状态码，非 1 状态码则为错误
	Meta    *Meta       `json:"meta,omitempty"` // 数据列表多了一个 Meta 字段
}

func (r Response) Decode(dest interface{}) (err error) {
	if !util.IsPoint(dest) {
		err = errors.New("decode: dest expect a point")
		return
	}

	var b []byte
	if b, err = json.Marshal(r.Data); err != nil {
		return
	}

	if err = json.Unmarshal(b, dest); err != nil {
		return
	}

	return
}
