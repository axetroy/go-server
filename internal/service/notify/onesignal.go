// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notify

// 推送通知模块
// 使用: onesignal 为推送中心

import (
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/onesignal"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
)

var sdk = onesignal.NewOneSignalClient(config.Notify.OneSignalAppID, config.Notify.OneSignalRestApiKey)

type Segment string                // 用户细分群体
type NotificationClickEvent string // 推送的点击事件

const (
	SegmentSubscribedUsers Segment = "Subscribed Users" // 所有已订阅的用户
	SegmentActiveUsers     Segment = "Active Users"     // 最近一周活跃的用户
	SegmentEngagedUsers    Segment = "Engaged Users"    // 最近一周重度依赖的用户
	SegmentInactiveUsers   Segment = "Inactive Users"   // 超过一周没有活跃的用户

	NotificationClickEventNone                  NotificationClickEvent = "none"                    // 空事件，点击通知什么都不会发送
	NotificationClickEventLoginAbnormal         NotificationClickEvent = "login_abnormal"          // 新的系统通知事件
	NotificationClickEventNewSystemNotification NotificationClickEvent = "new_system_notification" // 新的系统通知事件
	NotificationClickEventNewUserMessage        NotificationClickEvent = "new_user_message"        // 新的系统通知事件
)

// 发送推送附带的数据体结构
// event 给 APP 识别
// payload 是附带的数据
type NotificationBody struct {
	Event   NotificationClickEvent `json:"event"`   // 事件名
	Payload interface{}            `json:"payload"` // 数据体
}

func NewNotifierOneSignal() *NotifierOneSignal {
	n := NotifierOneSignal{}

	return &n
}

type NotifierOneSignal struct {
}

func (n *NotifierOneSignal) SendNotifyToAllUser(headings string, content string, data map[string]interface{}) error {
	err := sdk.CreateNotification(onesignal.CreateNotificationParams{
		IncludedSegments: []string{string(SegmentSubscribedUsers)},
		Headings:         map[string]string{"en": headings},
		Contents:         map[string]string{"en": content},
		Data: NotificationBody{
			Event:   NotificationClickEventNone,
			Payload: data,
		},
	})

	if err != nil {
		err = exception.ThirdParty.New(err.Error())
		return err
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifyToCustomUser(userId []string, headings string, content string, data map[string]interface{}) error {
	err := sdk.CreateNotification(onesignal.CreateNotificationParams{
		IncludeExternalUserIds: userId,
		Headings:               map[string]string{"en": headings},
		Contents:               map[string]string{"en": content},
		Data: NotificationBody{
			Event:   NotificationClickEventNone,
			Payload: data,
		},
	})

	if err != nil {
		err = exception.ThirdParty.New(err.Error())
		return err
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifySystemNotificationToUser(notificationId string) error {
	notificationInfo := model.Notification{}

	if err := database.Db.Model(notificationInfo).Where("id = ?", notificationId).First(&notificationInfo).Error; err != nil {
		// 如果没有这条系统通知，则跳过
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	err := sdk.CreateNotification(onesignal.CreateNotificationParams{
		IncludedSegments: []string{string(SegmentSubscribedUsers)},
		Headings:         map[string]string{"en": notificationInfo.Title},
		Contents:         map[string]string{"en": notificationInfo.Content},
		Data: NotificationBody{
			Event: NotificationClickEventNewSystemNotification,
			Payload: map[string]interface{}{
				"id":      notificationInfo.Id,
				"title":   notificationInfo.Title,
				"content": notificationInfo.Content,
			},
		},
	})

	if err != nil {
		err = exception.ThirdParty.New(err.Error())
		return err
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifyUserNewMessage(messageId string) error {
	messageInfo := model.Message{}

	if err := database.Db.Model(messageInfo).Where("id = ?", messageInfo).First(&messageInfo).Error; err != nil {
		// 如果没有这条消息，则跳过
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	err := sdk.CreateNotification(onesignal.CreateNotificationParams{
		IncludedSegments: []string{string(SegmentSubscribedUsers)},
		Headings:         map[string]string{"en": messageInfo.Title},
		Contents:         map[string]string{"en": messageInfo.Content},
		Data: NotificationBody{
			Event: NotificationClickEventNewUserMessage,
			Payload: map[string]interface{}{
				"id":      messageInfo.Id,
				"title":   messageInfo.Title,
				"content": messageInfo.Content,
			},
		},
	})

	if err != nil {
		err = exception.ThirdParty.New(err.Error())
		return err
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifyToUserForLoginStatus(userID string) error {
	var name string

	var userInfo = model.User{}

	if err := database.Db.Model(userInfo).Where("id = ?", userID).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果找不到用户，我们就跳过本次任务
			return nil
		}
		return err
	}

	if userInfo.Nickname == nil {
		name = userInfo.Username
	} else {
		name = *userInfo.Nickname
	}

	loginLogs := make([]model.LoginLog, 0)

	// 查找用户过往的登录记录, 只查找最近的两条
	if err := database.Db.Model(model.LoginLog{}).Where("uid = ?", userInfo.Id).Limit(2).Order("created_at DESC").Find(&loginLogs).Error; err != nil {
		// 如果没有之前的登录记录
		// 那么跳过
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}

	fmt.Printf("%+v\n", loginLogs)

	// 如果没有两条记录，那么不用作比较
	if len(loginLogs) < 2 {
		return nil
	}

	// 检查两次登录 IP 是否不一样
	// 如果两次 IP 一致，那么没有异常，跳过本次检查
	if loginLogs[0].Id == loginLogs[1].Id {
		return nil
	}

	err := sdk.CreateNotification(onesignal.CreateNotificationParams{
		IncludeExternalUserIds: []string{userInfo.Id},
		Headings:               map[string]string{"en": "异地登录异常"},
		Contents:               map[string]string{"en": fmt.Sprintf("发现您的帐号 [%s] 最近的登录异常，请注意帐号安全️", name)},
		Data: NotificationBody{
			Event:   NotificationClickEventLoginAbnormal,
			Payload: nil,
		},
	})

	if err != nil {
		err = exception.ThirdParty.New(err.Error())
		return err
	}

	fmt.Println("推送成功")

	return nil
}
