package admin_test

import (
	"fmt"
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogin(t *testing.T) {
	// 登陆超级管理员-失败
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin123",
		})

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.InvalidAccountOrPassword.Error(), r.Message)
		assert.Nil(t, r.Data)
	}

	// 登陆超级管理员-成功
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := admin.SignInResponse{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			// 判断UID是否与用户一致
			//c.Uid
			fmt.Println(c)
		}
	}
}
