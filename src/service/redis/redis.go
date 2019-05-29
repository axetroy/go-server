package redis

import (
	"github.com/axetroy/go-server/src/service/dotenv"
	"github.com/go-redis/redis"
	"os"
)

var (
	Client               *redis.Client // 默认的redis存储
	ActivationCodeClient *redis.Client // 存储激活码的
	ResetCodeClient      *redis.Client // 存储重置密码的
)

type redisConfig struct {
	Server   string
	Port     string
	Password string
}

func init() {
	if err := dotenv.Load(); err != nil {
		return
	}

	config := redisConfig{
		Server:   os.Getenv("REDIS_SERVER"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	var (
		addr     = config.Server + ":" + config.Port
		password = config.Password
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
