// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package oauth2

import (
	"errors"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"net/http"
	"net/url"
	"time"

	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func redirectToClient(c *router.Context, user *goth.User) {
	var (
		err      error
		tx       *gorm.DB
		finalURL string
	)

	frontendURL := dotenv.Get("OAUTH_REDIRECT_URL")

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			c.StatusCode(http.StatusBadRequest)
			c.Response(nil, schema.Response{Message: err.Error()})
		} else {
			c.Redirect(http.StatusTemporaryRedirect, finalURL)
		}
	}()

	uri, err := url.Parse(frontendURL)

	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		c.Response(nil, schema.Response{Message: "Invalid callback url"})
		return
	}

	tx = database.Db.Begin()

	oAuthInfo := model.OAuth{UserID: user.UserID}
	userInfo := model.User{}

	err = tx.Where(&oAuthInfo).First(&oAuthInfo).Error

	if err != nil {
		// 如果没找到对应的记录，说明这个帐号还没有绑定，我们给他创建一个本平台的帐号
		if err == gorm.ErrRecordNotFound {
			userName := user.Name

			if userName == "" {
				userName = user.NickName
			}

			if userName == "" {
				userName = user.FirstName + user.LastName
			}

			if userName == "" {
				userName = user.Provider + util.GenerateId()
			}

			userInfo = model.User{
				Username:                userName,
				Nickname:                &user.NickName,
				Password:                util.GeneratePassword(util.GenerateId()),
				Email:                   nil,
				Phone:                   nil,
				Status:                  model.UserStatusInit,
				UsernameRenameRemaining: 1,
			}

			// 创建一个用户
			if err = auth.CreateUserTx(tx, &userInfo, nil); err != nil {
				return
			}

			oAuthInfo.Uid = userInfo.Id
			oAuthInfo.Provider = model.OAuthProvider(user.Provider)
			oAuthInfo.Name = user.Name
			oAuthInfo.Nickname = user.NickName
			oAuthInfo.FirstName = user.FirstName
			oAuthInfo.LastName = user.LastName
			oAuthInfo.Description = user.Description
			oAuthInfo.Email = user.Email
			oAuthInfo.AvatarURL = user.AvatarURL
			oAuthInfo.Location = user.Location
			oAuthInfo.AccessToken = user.AccessToken
			oAuthInfo.AccessTokenSecret = user.AccessTokenSecret
			oAuthInfo.RefreshToken = user.RefreshToken
			oAuthInfo.ExpiresAt = user.ExpiresAt

			if err = tx.Create(&oAuthInfo).Error; err != nil {
				return
			}

		} else {
			return
		}
	}

	// 如果已经绑定帐号，则去查找帐号的相关信息
	if oAuthInfo.Uid != "" && userInfo.Id == "" {
		if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
			return
		}
	} else {
		// 如果有这条 oAuth 记录，但是没有这条绑定，这创建这个帐号
		userName := user.Name

		if userName == "" {
			userName = user.NickName
		}

		if userName == "" {
			userName = user.FirstName + user.LastName
		}

		if userName == "" {
			userName = user.Provider + util.GenerateId()
		}

		userInfo = model.User{
			Username:                userName,
			Nickname:                &user.NickName,
			Password:                util.GeneratePassword(util.GenerateId()),
			Email:                   nil,
			Phone:                   nil,
			Status:                  model.UserStatusInit,
			UsernameRenameRemaining: 1,
		}

		// 创建一个用户
		if err = auth.CreateUserTx(tx, &userInfo, nil); err != nil {
			return
		}
	}

	hash := util.MD5(user.UserID)

	if err := redis.ClientOAuthCode.Set(hash, userInfo.Id, time.Minute*5).Err(); err != nil {
		c.StatusCode(http.StatusBadRequest)
		c.Response(nil, schema.Response{
			Message: "Invalid callback url",
		})
		return
	}

	uri.Query().Set("access_token", hash)

	finalURL = uri.String()
}

var AuthRouter = router.Handler(func(c router.Context) {
	provider := c.Param("provider")

	c.ResetRequest(mux.SetURLVars(c.Request(), map[string]string{"provider": provider}))
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(c.Writer(), c.Request()); err == nil {
		// 认证成功
		redirectToClient(&c, &gothUser)
	} else {
		gothic.BeginAuthHandler(c.Writer(), c.Request())
	}
})

var AuthCallbackRouter = router.Handler(func(c router.Context) {
	provider := c.Param("provider")
	c.ResetRequest(mux.SetURLVars(c.Request(), map[string]string{"provider": provider}))
	gothUser, err := gothic.CompleteUserAuth(c.Writer(), c.Request())

	if err != nil {
		return
	}

	redirectToClient(&c, &gothUser)
})
