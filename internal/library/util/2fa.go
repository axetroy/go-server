// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package util

import (
	"crypto"
	"github.com/sec51/twofactor"
)

var (
	issuer     = "go-server" // 签发者
	encryption = crypto.SHA1 // 加密算法
	digits     = 6           // 密码位数
	prefix     = "prefix"    // 用于UID的前缀, 不能暴露这个字段，否则用户私钥可能泄漏
	suffix     = "suffix"    // 用户UID的后缀，不能暴露这个字段，否则用户私钥可能泄漏
)

// 拼接用户账号
func generateAccount(uid string) string {
	return prefix + uid + suffix
}

// 生成用户的密钥
func Generate2FASecret(uid string) (secret string, err error) {
	otp, err := twofactor.NewTOTP(generateAccount(uid), issuer, encryption, digits)

	if err != nil {
		return "", err
	}

	return otp.Secret(), nil
}

// 验证用户token是否正确
func Verify2FA(uid string, token string) bool {
	otp, err := twofactor.NewTOTP(generateAccount(uid), issuer, encryption, digits)

	if err != nil {
		return false
	}

	err = otp.Validate(token)

	return err == nil
}
