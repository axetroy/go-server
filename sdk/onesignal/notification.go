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

type Button struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Icon string `json:"icon,omitempty"`
	Url  string `json:"url,omitempty"`
}

type AndroidBackgroundLayout struct {
	Image         string `json:"image,omitempty"`
	HeadingsColor string `json:"headings_color,omitempty"`
	ContentsColor string `json:"contents_color,omitempty"`
}

// document: https://documentation.onesignal.com/reference/create-notification
type CreateNotificationParams struct {
	// Send to Segments
	IncludedSegments []string `json:"included_segments,omitempty"`
	ExcludedSegments []string `json:"excluded_segments,omitempty"`

	// Send to Specific Devices
	IncludePlayerIds       []string `json:"include_player_ids,omitempty"`
	IncludeExternalUserIds []string `json:"include_external_user_ids,omitempty"`
	IncludeEmailTokens     []string `json:"include_email_tokens,omitempty"`
	IncludeIosTokens       []string `json:"include_ios_tokens,omitempty"`
	IncludeWpWnsUris       []string `json:"include_wp_wns_uris,omitempty"`
	IncludeAmazonRegIds    []string `json:"include_amazon_reg_ids,omitempty"`
	IncludeChromeRegIds    []string `json:"include_chrome_reg_ids,omitempty"`
	IncludeChromeWebRegIds []string `json:"include_chrome_web_reg_ids,omitempty"`
	IncludeAndroidRegIds   []string `json:"include_android_reg_ids,omitempty"`
	ExternalId             string   `json:"external_id,omitempty"`

	// Formatting Filters
	Filters []Filter `json:"filters,omitempty"`

	// Content & Language
	Contents         map[string]string `json:"contents,omitempty"`
	Headings         map[string]string `json:"headings,omitempty"`
	Subtitle         map[string]string `json:"subtitle,omitempty"`
	TemplateId       string            `json:"template_id,omitempty"`
	ContentAvailable bool              `json:"content_available,omitempty"`
	MutableContent   bool              `json:"mutable_content,omitempty"`

	// Email Content
	EmailSubject     string `json:"email_subject,omitempty"`
	EmailBody        string `json:"email_body,omitempty"`
	EmailFromName    string `json:"email_from_name,omitempty"`
	EmailFromAddress string `json:"email_from_address,omitempty"`

	// Attachments
	Data             map[string]interface{} `json:"data,omitempty"`
	Url              string                 `json:"url,omitempty"`
	WebUrl           string                 `json:"web_url,omitempty"`
	AppUrl           string                 `json:"app_url,omitempty"`
	IOSAttachments   string                 `json:"ios_attachments,omitempty"`
	BigPicture       string                 `json:"big_picture,omitempty"`
	ChromeWebImage   string                 `json:"chrome_web_image,omitempty"`
	AdmBigPicture    string                 `json:"adm_big_picture,omitempty"`
	ChromeBigPicture string                 `json:"chrome_big_picture,omitempty"`

	// Action Buttons
	Buttons     []Button `json:"buttons,omitempty"`
	WebButtons  []Button `json:"web_buttons,omitempty"`
	IOSCategory string   `json:"ios_category,omitempty"`

	// Appearance
	AndroidChannelID         string                  `json:"android_channel_id,omitempty"`
	ExistingAndroidChannelID string                  `json:"existing_android_channel_id,omitempty"`
	AndroidBackgroundLayout  AndroidBackgroundLayout `json:"android_background_layout,omitempty"`
	SmallIcon                string                  `json:"small_icon,omitempty"`
	LargeIcon                string                  `json:"large_icon,omitempty"`
	AdmSmallIcon             string                  `json:"adm_small_icon,omitempty"`
	AdmLargeIcon             string                  `json:"adm_large_icon,omitempty"`
	ChromeWebIcon            string                  `json:"chrome_web_icon,omitempty"`
	//ChromeWebImage           string                  `json:"chrome_web_image,omitempty"`
	ChromeWebBadge     string `json:"chrome_web_badge,omitempty"`
	FirefoxIcon        string `json:"firefox_icon,omitempty"`
	ChromeIcon         string `json:"chrome_icon,omitempty"`
	IOSSound           string `json:"ios_sound,omitempty"`
	AndroidSound       string `json:"android_sound,omitempty"`
	AdmSound           string `json:"AdmSound,omitempty"`
	WpWnsSound         string `json:"wp_wns_sound,omitempty"`
	AndroidLedColor    string `json:"android_led_color,omitempty"`
	AndroidAccentColor string `json:"android_accent_color,omitempty"`
	AndroidVisibility  string `json:"android_visibility,omitempty"`
	IOSBadgeType       string `json:"ios_badgeType,omitempty"`
	IOSBadgeCount      string `json:"ios_badgeCount,omitempty"`
	CollapseID         string `json:"collapse_id,omitempty"`
	ApnsAlert          string `json:"apns_alert,omitempty"`

	// Delivery
	SendAfter            string `json:"send_after,omitempty"`
	DelayedOption        string `json:"delayed_option,omitempty"`
	DeliveryTimeOfDay    string `json:"delivery_time_of_day,omitempty"`
	TTL                  int    `json:"ttl,omitempty"`
	Priority             int    `json:"priority,omitempty"`
	ApnsPushTypeOverride string `json:"apns_push_type_override,omitempty"`

	// Grouping & Collapsing
	AndroidGroup        string            `json:"android_group,omitempty"`
	AndroidGroupMessage string            `json:"android_group_message,omitempty"`
	AdmGroup            string            `json:"adm_group,omitempty"`
	AdmGroupMessage     map[string]string `json:"adm_group_message,omitempty"`
	ThreadID            string            `json:"thread_id,omitempty"`
	SummaryArg          string            `json:"summary_arg,omitempty"`
	SummaryArgCount     float64           `json:"summary_arg_count,omitempty"`

	// Platform to Deliver To
	IsIOS                     bool   `json:"isIos,omitempty"`
	IsAndroid                 bool   `json:"isAndroid,omitempty"`
	IsAnyWeb                  bool   `json:"isAnyWeb,omitempty"`
	IsEmail                   bool   `json:"isEmail,omitempty"`
	IsChromeWeb               bool   `json:"isChromeWeb,omitempty"`
	IsFirefox                 bool   `json:"isFirefox,omitempty"`
	IsSafari                  bool   `json:"isSafari,omitempty"`
	IsWP_WNS                  bool   `json:"isWP_WNS,omitempty"`
	IsAdm                     bool   `json:"isAdm,omitempty"`
	IsChrome                  bool   `json:"isChrome,omitempty"`
	ChannelForExternalUserIds string `json:"channel_for_external_user_ids,omitempty"`
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
