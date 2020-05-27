// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package onesignal

func NewOneSignalClient(AppID string, RestApiKey string) *OneSignal {
	return &OneSignal{
		appId:      AppID,
		restApiKey: RestApiKey,
	}
}

type OneSignal struct {
	appId      string
	restApiKey string
}
