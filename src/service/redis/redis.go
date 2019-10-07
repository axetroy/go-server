// Copyright 2019 Axetroy. All rights reserved. MIT license.
package redis

import (
	"github.com/axetroy/go-server/src/config"
	"github.com/go-redis/redis"
)

var (
	Client               *redis.Client // 默认的redis存储
	ClientActivationCode *redis.Client // 存储帐号激活码的
	ClientResetCode      *redis.Client // 存储重置密码的
	Config               = config.Redis
)

func init() {
	var (
		addr     = Config.Host + ":" + Config.Port
		password = Config.Password
	)

	// 初始化3个DB连接
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

}
