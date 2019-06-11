// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/email"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGenerateResetCode(t *testing.T) {
	user, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(user.Username)

	code := email.GenerateResetCode(user.Id)

	assert.IsType(t, "", code)
}

func TestSendResetPasswordEmail(t *testing.T) {

	body, _ := json.Marshal(&email.SendResetPasswordEmailParams{
		To: "123adsd@dasdad.com", // invalid email
	})

	r := tester.HttpUser.Post("/v1/email/send/password/reset", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, common_error.ErrUserNotExist.Error(), res.Message)
}
