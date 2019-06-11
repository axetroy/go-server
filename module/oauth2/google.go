// Copyright 2019 Axetroy. All rights reserved. MIT license.
package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/config"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

var googleOAuth2Config *oauth2.Config

func GetgoogleOAuth2Config() oauth2.Config {
	var endpoint = oauth2.Endpoint{
		AuthURL:   "https://accounts.google.com/o/oauth2/auth",
		TokenURL:  "https://oauth2.googleapis.com/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}
	if googleOAuth2Config != nil {
		return *googleOAuth2Config
	}

	googleOAuth2Config = &oauth2.Config{
		ClientID:     config.OAuth2Google.ClientId,
		ClientSecret: config.OAuth2Google.ClientSecret,
		RedirectURL:  config.User.Domain + "/v1/oauth2/google_callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: endpoint,
	}

	return *googleOAuth2Config
}

const oauthStateString = "go-server"

// 调用谷歌登陆，然后重定向到谷歌认证页面
func GoogleLoginRouter(ctx *gin.Context) {
	c := GetgoogleOAuth2Config()
	url := c.AuthCodeURL(oauthStateString)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

type Query struct {
	State string `form:"state" json:"state"`
	Code  string `form:"code" json:"code"`
}

type GoogleAuthResponse struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

// 谷歌登陆成功之后的回调函数
func GoogleCallbackRouter(ctx *gin.Context) {
	query := Query{}

	if err := ctx.BindQuery(&query); err != nil {
		fmt.Printf("error")
		return
	}

	if query.State != oauthStateString {
		res := fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, query.State)
		ctx.String(http.StatusBadRequest, res)
		return
	}

	c := GetgoogleOAuth2Config()

	token, err := c.Exchange(oauth2.NoContext, query.Code)

	if err != nil {
		res := fmt.Sprintf("code exchange failed with '%s'\n", err)
		ctx.String(http.StatusTemporaryRedirect, res)
		return
	}

	// 在中国有防火墙，访问不了Google
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	res := &GoogleAuthResponse{}

	err = json.Unmarshal(contents, &res)

	if err != nil {
		ctx.String(http.StatusTemporaryRedirect, err.Error())
	}

	// 查询是否有这个用户存在
	user := user_model.User{OauthGoogleId: &res.Id}

	if err = database.Db.Where(&user).Last(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果用户不存在，则创建一个用户
			userInfo := user_model.User{
				Username: res.Name,
				Nickname: &res.Name,
				Password: util.GeneratePassword(util.RandomString(8)), // 生成一个随机密码
				Status:   user_model.UserStatusInit,
				Email:    &res.Email,
				Gender:   user_model.GenderErrUnknown,
			}

			if err = database.Db.Create(&userInfo).Error; err != nil {
				return
			}
		}
		return
	}

	// TODO：重定向到前端页面
}
