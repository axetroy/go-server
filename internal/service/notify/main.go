// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notify

import "github.com/axetroy/go-server/internal/schema"

type Content struct {
	EN string `json:"en"`
}

type Headings struct {
	EN string `json:"en"`
}

type SendNotifyEvent string

const (
	SendNotifyEventSendNotifyToAllUser           SendNotifyEvent = "SendNotifyToAllUser"
	SendNotifyEventSendNotifyToCustomUser        SendNotifyEvent = "SendNotifyToCustomUser"
	SendNotifyEventSendNotifyToLoginAbnormalUser SendNotifyEvent = "SendNotifyToLoginAbnormalUser"
)

type Notifier interface {
	SendNotifyToAllUser(headings string, content string) error                      // 向所有用户推送
	SendNotifyToCustomUser(userIds []string, headings string, content string) error // 推送自定义通知
	//SendNotifySystemMessageToUser(userIds []string, title string, content string) error // 推送系统通知
	//SendNotifyUserMessageToUser(userIds []string, title string, content string) error   // 推送用户消息
	SendNotifyToLoginAbnormalUser(userInfo schema.ProfilePublic) error // 推送用户登录异常
}

var Notify *Notifier

func init() {
	Notify = getInstance(NewNotifierOneSignal())
}

func getInstance(n Notifier) *Notifier {
	return &n
}
