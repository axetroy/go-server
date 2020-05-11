// Copyright 2019 Axetroy. All rights reserved. MIT license.
package redis

import (
	"github.com/axetroy/go-server/internal/config"
	"github.com/go-redis/redis"
)

var (
	Client               *redis.Client // 默认的redis存储
	ClientActivationCode *redis.Client // 存储帐号激活码的
	ClientAuthEmailCode  *redis.Client // 存储邮箱验证码，存储结构 key: 验证码, value: 邮箱
	ClientAuthPhoneCode  *redis.Client // 存储手机验证码，存储结构 key: 验证码, value: 手机号
	ClientResetCode      *redis.Client // 存储重置密码的
	ClientOAuthCode      *redis.Client // 存储 oAuth2 对应的激活码
	Config               = config.Redis
)

func init() {
	var (
		addr     = Config.Host + ":" + Config.Port
		password = Config.Password
	)

	// 初始化DB连接
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	ClientActivationCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	ClientResetCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})

	ClientAuthEmailCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       3,
	})

	ClientAuthPhoneCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       4,
	})

	ClientOAuthCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       5,
	})

}
