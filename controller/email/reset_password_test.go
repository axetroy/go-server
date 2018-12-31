package email_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/axetroy/server/controller/email"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/tester"
	"net/http"
	"testing"
)

func TestGenerateResetCode(t *testing.T) {
	code := email.GenerateResetCode(tester.Uid)

	assert.IsType(t, "", code)
}

func TestSendResetPasswordEmail(t *testing.T) {

	body, _ := json.Marshal(&email.SendResetPasswordEmailParams{
		To: "123adsd@dasdad.com", // invalid email
	})

	r := tester.Http.Post("/v1/email/send/reset_password", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.UserNotExist.Error(), res.Message)
}
