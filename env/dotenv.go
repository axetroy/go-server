package env

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"path"
)

var (
	Test    bool
	RootDir string
)

func Load() (err error) {
	var envFilePath = ".env"
	isRunInTest := flag.Lookup("test.v") != nil

	Test = isRunInTest

	if isRunInTest {
		RootDir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "axetroy", "go-server")
		envFilePath = path.Join(RootDir, envFilePath)
	}

	err = godotenv.Load(envFilePath)
	return
}
