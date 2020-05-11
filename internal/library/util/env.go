package util

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/fatih/color"
	"log"
	"os"
	"runtime"
)

func init() {
	err := dotenv.Load()

	if err != nil {
		panic(err)
	}
}

func printJSON(o interface{}) {
	if output, err := json.Marshal(o); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(output))
	}
}

func PrintEnv() {
	envs := os.Environ()

	fmt.Println(color.GreenString("=== Runtime ==="))

	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Go OS: %s\n", runtime.GOOS)
	fmt.Printf("Go Arch: %s\n", runtime.GOARCH)

	fmt.Println(color.GreenString("=== Environmental Variable ==="))

	for _, e := range envs {
		fmt.Println(e)
	}

	fmt.Println(color.GreenString("=== Configuration Common ==="))
	printJSON(config.Common)

	fmt.Println(color.GreenString("=== Configuration Upload ==="))
	printJSON(config.Upload)

	fmt.Println(color.GreenString("=== Configuration Message Queue ==="))
	printJSON(config.MessageQueue)

	fmt.Println(color.GreenString("=== Configuration Database ==="))
	printJSON(config.Database)

	fmt.Println(color.GreenString("=== Configuration Wechat ==="))
	printJSON(config.Wechat)

	fmt.Println(color.GreenString("=== Configuration Redis ==="))
	printJSON(config.Redis)

	fmt.Println(color.GreenString("=== Configuration SMTP ==="))
	printJSON(config.SMTP)

	fmt.Println(color.GreenString("=== Configuration User ==="))
	printJSON(config.User)

	fmt.Println(color.GreenString("=== Configuration Admin ==="))
	printJSON(config.Admin)
}
