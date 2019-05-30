package redis

import (
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/service/dotenv"
	"github.com/go-redis/redis"
)

var (
	Client               *redis.Client // 默认的redis存储
	ActivationCodeClient *redis.Client // 存储激活码的
	ResetCodeClient      *redis.Client // 存储重置密码的
	Config               = config.Redis
)

func init() {
	if err := dotenv.Load(); err != nil {
		return
	}

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

	ActivationCodeClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	ResetCodeClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})

}
