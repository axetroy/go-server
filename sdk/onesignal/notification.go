// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package onesignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FilterField string

type Filter struct {
	Field    string `json:"field"`
	Key      string `json:"key"`
	Relation string `json:"relation"`
	Value    string `json:"value"`
}

type CreateNotificationParams struct {
	// optional
	IncludedSegments       []string               `json:"included_segments,omitempty"`
	ExcludedSegments       []string               `json:"excluded_segments,omitempty"`
	IncludePlayerIds       []string               `json:"include_player_ids,omitempty"`
	IncludeExternalUserIds []string               `json:"include_external_user_ids,omitempty"`
	IncludeEmailTokens     []string               `json:"include_email_tokens,omitempty"`
	IncludeIosTokens       []string               `json:"include_ios_tokens,omitempty"`
	IncludeWpWnsUris       []string               `json:"include_wp_wns_uris,omitempty"`
	IncludeAmazonRegIds    []string               `json:"include_amazon_reg_ids,omitempty"`
	IncludeChromeRegIds    []string               `json:"include_chrome_reg_ids,omitempty"`
	IncludeChromeWebRegIds []string               `json:"include_chrome_web_reg_ids,omitempty"`
	IncludeAndroidRegIds   []string               `json:"include_android_reg_ids,omitempty"`
	ExternalId             string                 `json:"external_id,omitempty"`
	Filters                []Filter               `json:"filters,omitempty"`
	Data                   map[string]interface{} `json:"data,omitempty"`

	// required
	Contents         map[string]string `json:"contents,omitempty"`
	Headings         map[string]string `json:"headings,omitempty"`
	Subtitle         map[string]string `json:"subtitle,omitempty"`
	TemplateId       string            `json:"template_id,omitempty"`
	ContentAvailable bool              `json:"content_available,omitempty"`
	MutableContent   bool              `json:"mutable_content,omitempty"`
}

type NotificationResponse struct {
	// see: https://documentation.onesignal.com/reference/create-notification
	ID         string      `json:"id"`
	Recipients int         `json:"recipients"` // 接收到的数量
	ExternalID *string     `json:"external_id"`
	Errors     interface{} `json:"errors"`
}

func (o *OneSignal) CreateNotification(params CreateNotificationParams) error {
	bodyByte, err := json.Marshal(params)

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewReader(bodyByte))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", o.restApiKey))

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

	resFromRemote := NotificationResponse{}

	if err := json.Unmarshal(resByte, &resFromRemote); err != nil {
		return err
	}

	if resFromRemote.Errors != nil {
		switch t := resFromRemote.Errors.(type) {
		case []string:
			return errors.New(t[0])
		case map[string]string:
			return errors.New("unknown error")
		}
	}

	return nil
}
