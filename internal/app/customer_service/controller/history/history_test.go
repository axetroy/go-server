// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package history_test

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/history"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 获取客服的会话记录
func TestGetWaiterSession(t *testing.T) {
	// 创建会话记录
	result, err := history.GetWaiterSession("278699413700870144")

	assert.Nil(t, err)

	fmt.Printf("%+v\n", result)
}

// 获取聊天记录
func TestGetHistory(t *testing.T) {

}
