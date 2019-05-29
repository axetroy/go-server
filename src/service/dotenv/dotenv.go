package dotenv

import (
	"flag"
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var (
	Test    bool   // 当前是否是测试环境
	Env     string // 当前的运行环境
	RootDir string // 当前运行的二进制所在的目录
	loaded  bool   // 是否已初始化过
)

func Load() (err error) {
	if loaded {
		return
	}
	defer func() {
		loaded = true
	}()
	isRunInTest := flag.Lookup("test.v") != nil
	isRunInTravis := os.Getenv("TRAVIS") != ""

	Test = isRunInTest
	Env = os.Getenv("GO_ENV")

	var pwd string

	pwd, err = os.Getwd()

	if err != nil {
		return
	}

	switch true {
	// 如果运行才 travis，则取当前目录
	case isRunInTravis:
		Test = true
		RootDir = os.Getenv("TRAVIS_BUILD_DIR")
		break
	// 如果运行在测试用例
	case isRunInTest:
		RootDir = pwd
		break
	default:
		ex, err := os.Executable()

		if err != nil {
			panic(err)
		}

		exPath := filepath.Dir(ex)

		RootDir = filepath.Dir(exPath)

		// 如果是以 go run main.go 运行, 则取工作目录
		goRunReg := regexp.MustCompile("go-build\\d+")
		// 如果是运行在 idea 里面的话
		ideaRunReg := regexp.MustCompile("___go_build_")

		isRunInGoRun := goRunReg.MatchString(ex)
		isRunInIdea := ideaRunReg.MatchString(ex)

		switch true {
		case isRunInGoRun:
		case isRunInIdea:
			RootDir = pwd
		}
	}

	fmt.Println("当前工作目录", RootDir)

	dotEnvFilePath := path.Join(RootDir, ".env")

	if !fs.PathExists(dotEnvFilePath) {
		return
	}

	err = godotenv.Load(dotEnvFilePath)
	return
}
