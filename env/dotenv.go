package env

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"path"
)

func Load() (err error) {
	var envFilePath = ".env"
	isRunInTest := flag.Lookup("test.v") != nil

	if isRunInTest {
		envFilePath = path.Join(os.Getenv("GOPATH"), "src", "github.com", "axetroy", "go-server", envFilePath)
	}

	err = godotenv.Load(envFilePath)
	return
}
