package util_test

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	testPassword := "password"
	s := util.GeneratePassword(testPassword)

	// 生成的密码与预期的一致
	if !assert.Equal(t, "39d9782aa70de6b4944c40991ac37004", s) {
		return
	}

	// 生成两次的密码保持一致
	if !assert.Equal(t, util.GeneratePassword(testPassword), util.GeneratePassword(testPassword)) {
		return
	}

	if !assert.Len(t, util.GeneratePassword(testPassword), 32, "密码必须是32位") {
		return
	}
}
