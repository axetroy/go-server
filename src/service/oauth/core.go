package oauth

import (
	"fmt"
	"github.com/axetroy/go-server/src/service/dotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitter"
)

const (
	domain       = "http://localhost:3000"
	callbackPath = "/v1/oauth2/%s/callback"
)

func generateCallbackURL(provider string) string {
	return domain + fmt.Sprintf(callbackPath, provider)
}

func init() {

	goth.UseProviders(
		twitter.New(dotenv.Get("TWITTER_KEY"), dotenv.Get("TWITTER_SECRET"), generateCallbackURL("twitter")),
		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(dotenv.Get("TWITTER_KEY"), dotenv.Get("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),

		facebook.New(dotenv.Get("FACEBOOK_KEY"), dotenv.Get("FACEBOOK_SECRET"), generateCallbackURL("facebook")),
		google.New(dotenv.Get("GOOGLE_KEY"), dotenv.Get("GOOGLE_SECRET"), generateCallbackURL("google")),
		github.New("GITHUB_KEY", "GITHUB_SECRET", generateCallbackURL("github")),
		gitlab.New(dotenv.Get("GITLAB_KEY"), dotenv.Get("GITLAB_SECRET"), generateCallbackURL("gitlab")),
	)

	// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
	// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
	// ignore the error for now
	//connection, _ := openidConnect.New(dotenv.Get("OPENID_CONNECT_KEY"), dotenv.Get("OPENID_CONNECT_SECRET"), "http://localhost:3000/auth/openid-connect/callback", dotenv.Get("OPENID_CONNECT_DISCOVERY_URL"))
	//
	//if connection != nil {
	//	goth.UseProviders(connection)
	//}
}
