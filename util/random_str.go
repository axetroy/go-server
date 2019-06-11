// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util

import (
	"math/rand"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano()) // 不设置随机种子的话，随机数都是一样的
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
