// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notify

// 推送通知模块
// 使用: onesignal 为推送中心

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"io/ioutil"
	"net/http"
)

var (
	appId      = config.Notify.OneSignalAppID
	restApiKey = config.Notify.OneSignalRestApiKey
)

type oneSignalResponse struct {
	// see: https://documentation.onesignal.com/reference/create-notification
	ID         string      `json:"id"`
	Recipients int         `json:"recipients"` // 接收到的数量
	ExternalID *string     `json:"external_id"`
	Errors     interface{} `json:"errors"`
}

func NewNotifierOneSignal() *NotifierOneSignal {
	n := NotifierOneSignal{}

	return &n
}

type NotifierOneSignal struct {
}

func (n *NotifierOneSignal) SendNotifyToAllUser(headings string, content string) error {
	type Body struct {
		AppID    string   `json:"app_id"`
		Contents Content  `json:"contents"`
		Headings Headings `json:"headings"`
	}

	var body = Body{
		AppID:    appId,
		Contents: Content{EN: content},
		Headings: Headings{EN: headings},
	}

	bodyByte, err := json.Marshal(body)

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewReader(bodyByte))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", restApiKey))

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	resByte, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		msg := resByte

		if len(msg) == 0 {
			msg = []byte(http.StatusText(res.StatusCode))
		}

		return errors.New(string(msg))
	}

	resFromRemote := oneSignalResponse{}

	if err := json.Unmarshal(resByte, &resFromRemote); err != nil {
		return err
	}

	if resFromRemote.Errors != nil {
		switch t := resFromRemote.Errors.(type) {
		case []string:
			return errors.New(t[0])
		case map[string]string:
			return exception.Unknown
		}
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifyToCustomUser(userId []string, headings string, content string) error {
	type Body struct {
		IncludeExternalUserIds []string `json:"include_external_user_ids"`
		AppID                  string   `json:"app_id"`
		Contents               Content  `json:"contents"`
		Headings               Headings `json:"headings"`
	}

	var body = Body{
		IncludeExternalUserIds: userId,
		AppID:                  appId,
		Contents:               Content{EN: content},
		Headings:               Headings{EN: headings},
	}

	bodyByte, err := json.Marshal(body)

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewReader(bodyByte))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", restApiKey))

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	resByte, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		msg := resByte

		if len(msg) == 0 {
			msg = []byte(http.StatusText(res.StatusCode))
		}

		return errors.New(string(msg))
	}

	resFromRemote := oneSignalResponse{}

	if err := json.Unmarshal(resByte, &resFromRemote); err != nil {
		return err
	}

	if resFromRemote.Errors != nil {
		switch t := resFromRemote.Errors.(type) {
		case []string:
			return errors.New(t[0])
		case map[string]string:
			return exception.Unknown
		}
	}

	return nil
}

func (n *NotifierOneSignal) SendNotifyToLoginAbnormalUser(userInfo schema.ProfilePublic) error {
	type Body struct {
		IncludeExternalUserIds []string `json:"include_external_user_ids"`
		AppID                  string   `json:"app_id"`
		Contents               Content  `json:"contents"`
		Headings               Headings `json:"headings"`
	}

	var name string

	if userInfo.Nickname == nil {
		name = userInfo.Username
	} else {
		name = *userInfo.Nickname
	}

	var body = Body{
		IncludeExternalUserIds: []string{userInfo.Id},
		AppID:                  appId,
		Contents: Content{
			EN: fmt.Sprintf("发现您的帐号 [%s] 最近的登录异常，请注意帐号安全️", name),
		},
		Headings: Headings{
			EN: "异地登录异常⚠",
		},
	}

	bodyByte, err := json.Marshal(body)

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewReader(bodyByte))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", restApiKey))

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	resByte, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		msg := resByte

		if len(msg) == 0 {
			msg = []byte(http.StatusText(res.StatusCode))
		}

		return errors.New(string(msg))
	}

	resFromRemote := oneSignalResponse{}

	if err := json.Unmarshal(resByte, &resFromRemote); err != nil {
		return err
	}

	if resFromRemote.Errors != nil {
		switch t := resFromRemote.Errors.(type) {
		case []string:
			return errors.New(t[0])
		case map[string]string:
			return exception.Unknown
		}
	}

	return nil
}
