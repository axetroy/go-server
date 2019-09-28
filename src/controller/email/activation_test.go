// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/email"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGenerateActivationCode(t *testing.T) {
	user, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(user.Username)

	code := email.GenerateResetCode(user.Id)

	assert.IsType(t, "", code)
}

func TestSendActivationEmail(t *testing.T) {

	body, _ := json.Marshal(&email.SendActivationEmailParams{
		To: "123adsd@dasdad.com", // invalid email
	})

	r := tester.HttpUser.Post("/v1/email/send/activation", body, nil)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
		return
	}

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.UserNotExist.Error(), res.Message)
}
