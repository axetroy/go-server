package email_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/email"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGenerateActivationCode(t *testing.T) {
	testerUsername := "tester-TestGenerateActivationCode"
	testerUid := ""

	// 动态创建一个测试账号
	{
		r := auth.SignUp(auth.SignUpParams{
			Username: &testerUsername,
			Password: "123123",
		})

		profile := user.Profile{}

		assert.Nil(t, tester.Decode(r.Data, &profile))

		testerUid = profile.Id

		defer func() {
			auth.DeleteUserByUserName(testerUsername)
		}()
	}

	code := email.GenerateResetCode(testerUid)

	assert.IsType(t, "", code)
}

func TestSendActivationEmail(t *testing.T) {

	body, _ := json.Marshal(&email.SendActivationEmailParams{
		To: "123adsd@dasdad.com", // invalid email
	})

	r := tester.Http.Post("/v1/email/send/activation", body, nil)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := response.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, response.StatusFail, res.Status) {
		fmt.Println(res.Message)
		return
	}
	if !assert.Equal(t, exception.UserNotExist.Error(), res.Message) {
		return
	}
}
