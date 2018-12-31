package redis

import (
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"os"
)

var (
	Common         *redis.Client // 默认的redis存储
	ActivationCode *redis.Client // 存储激活码的
	ResetCode      *redis.Client // 存储重置密码的
)

type Config struct {
	Server   string
	Port     string
	Password string
}

var config Config

func init() {

	if err := godotenv.Load(); err != nil {
		return
	}

	config = Config{
		Server:   os.Getenv("REDIS_SERVER"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	var (
		addr     = config.Server + ":" + config.Port
		password = config.Password
	)

	println(addr, password)

	Common = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	ActivationCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	ResetCode = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})

}
