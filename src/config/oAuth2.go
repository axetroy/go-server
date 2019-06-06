// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
)

type oAuth2Google struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var OAuth2Google oAuth2Google

func init() {
	OAuth2Google.ClientId = dotenv.Get("GOOGLE_AUTH2_CLIENT_ID")
	OAuth2Google.ClientSecret = dotenv.Get("GOOGLE_AUTH2_CLIENT_SECRET")
}
