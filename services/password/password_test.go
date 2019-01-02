package password_test

import (
	"fmt"
	"github.com/axetroy/go-server/services/password"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	fmt.Println("运行测试用例")
	testPassword := "password"
	s := password.Generate(testPassword)

	// 生成的密码与预期的一致
	if !assert.Equal(t, "39d9782aa70de6b4944c40991ac37004", s) {
		return
	}

	// 生成两次的密码保持一致
	if !assert.Equal(t, password.Generate(testPassword), password.Generate(testPassword)) {
		return
	}

	if !assert.Len(t, password.Generate(testPassword), 32, "密码必须是32位") {
		return
	}

	fmt.Println(password.Generate("123123"))
}
