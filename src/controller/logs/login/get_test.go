package login_test

import (
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/logs/login"
	"github.com/axetroy/go-server/tester"
	"testing"
)

func TestGetLoginLog(t *testing.T) {
	user, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(user.Username)

	r := login.GetLatestLog(&controller.Context{Uid: user.Id})

	fmt.Printf("%+v\n", r)
}
