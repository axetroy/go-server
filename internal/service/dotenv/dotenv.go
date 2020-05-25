// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package dotenv

import (
	"flag"
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	Test    bool   // 当前是否是测试环境
	RootDir string // 当前运行的二进制所在的目录
	loaded  bool   // 是否已初始化过
)

func init() {
	if err := Load(); err != nil {
		panic(err)
	}

	if !Test {
		if os.Getenv("GO_TESTING") != "" || strings.Index(os.Getenv("XPC_SERVICE_NAME"), "com.jetbrains.goland") >= 0 {
			Test = true
		}
	}
}

func Load() (err error) {
	if loaded {
		return
	}
	defer func() {
		if err == nil {
			loaded = true
		}
	}()
	isRunInTest := flag.Lookup("test.v") != nil
	isRunInTravis := os.Getenv("TRAVIS") != ""

	// 如果设置环境变量 GO_TESTING=1 则认为是测试环境
	if !isRunInTest {
		isRunInTest = len(os.Getenv("GO_TESTING")) > 0
	}

	if !isRunInTest {
		if isRunInTravis {
			isRunInTest = true
		} else {
			e, _ := os.Executable()
			isRunInTest = regexp.MustCompile("\\/T\\/___").MatchString(e)
		}
	}

	Test = isRunInTest

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

		exPathDir := filepath.Dir(ex)

		RootDir = exPathDir

		// 如果是以 go run main.go 运行, 则取工作目录
		goRunReg := regexp.MustCompile("/go-build\\d+/")
		// 如果是运行在 IDEA 里面的话
		ideaRunReg := regexp.MustCompile("___go_build_")

		ifRunInGoRun := goRunReg.MatchString(ex)
		ifRunInIdea := ideaRunReg.MatchString(ex)

		switch true {
		case ifRunInGoRun:
			RootDir = pwd
			break
		case ifRunInIdea:
			RootDir = pwd
			break
		}
	}

	dotEnvFilePath := path.Join(RootDir, ".env")

	if !fs.PathExists(dotEnvFilePath) {
		return
	}

	fmt.Println(fmt.Sprintf("加载环境变量文件 `%s`", color.GreenString(dotEnvFilePath)))

	err = godotenv.Load(dotEnvFilePath)
	return
}

func Get(key string) string {
	if loaded == false {
		_ = Load()
	}
	return os.Getenv(key)
}

func GetIntByDefault(key string, defaultValue int) int {
	val := GetByDefault(key, fmt.Sprintf("%d", defaultValue))

	result, err := strconv.Atoi(val)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetInt64ByDefault(key string, defaultValue int64) int64 {
	val := GetByDefault(key, fmt.Sprintf("%d", defaultValue))

	result, err := strconv.Atoi(val)

	if err != nil {
		log.Fatal(err)
	}

	return int64(result)
}

func GetStrArrayByDefault(key string, defaultValue []string) []string {
	val := GetByDefault(key, fmt.Sprintf("%s", strings.Join(defaultValue, ",")))

	var result []string

	arr := strings.Split(val, ",")

	for _, val := range arr {
		result = append(result, strings.TrimSpace(val))
	}

	return result
}

func GetByDefault(key string, defaultValue string) string {
	if loaded == false {
		_ = Load()
	}
	result := os.Getenv(key)

	if result == "" {
		return defaultValue
	} else {
		return result
	}
}
