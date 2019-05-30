package config

import "os"

type oAuth2Google struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var OAuth2Google oAuth2Google

func init() {
	OAuth2Google.ClientId = os.Getenv("GOOGLE_AUTH2_CLIENT_ID")
	OAuth2Google.ClientSecret = os.Getenv("GOOGLE_AUTH2_CLIENT_SECRET")
}
