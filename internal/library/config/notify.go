// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type notify struct {
	// https://documentation.onesignal.com/reference/create-notification
	OneSignalAppID      string `json:"one_signal_app_id"`
	OneSignalRestApiKey string `json:"one_signal_rest_api_key"`
}

var Notify notify

func init() {
	Notify.OneSignalAppID = dotenv.GetByDefault("ONE_SIGNAL_APP_ID", "")
	Notify.OneSignalRestApiKey = dotenv.GetByDefault("ONE_SIGNAL_REST_API_KEY", "")
}
