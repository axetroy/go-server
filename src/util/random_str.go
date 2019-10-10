// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers     = "0123456789"
)

func RandomString(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano()) // 不设置随机种子的话，随机数都是一样的
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomNumeric(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano()) // 不设置随机种子的话，随机数都是一样的
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}
