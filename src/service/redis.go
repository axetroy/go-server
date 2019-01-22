package service

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/go-redis/redis"
	"os"
)

var (
	RedisClient               *redis.Client // 默认的redis存储
	RedisActivationCodeClient *redis.Client // 存储激活码的
	RedisResetCodeClient      *redis.Client // 存储重置密码的
)

type redisConfig struct {
	Server   string
	Port     string
	Password string
}

func init() {
	if err := util.LoadEnv(); err != nil {
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
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	RedisActivationCodeClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	RedisResetCodeClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})

}
