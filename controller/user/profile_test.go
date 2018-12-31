package user_test

import (
	"encoding/json"
	"github.com/axetroy/mocker"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gitlab.com/axetroy/server/controller/user"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/tester"
	"net/http"
	"testing"
)

func TestGetProfileWithInvalidAuth(t *testing.T) {

	header := mocker.Header{
		"Authorization": "Bearera 12312", // invalid Bearera
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidAuth.Error(), res.Message)
}

func TestGetProfileWithInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": "Bearer 12312",
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidToken.Error(), res.Message)
}

func TestGetProfile(t *testing.T) {
	header := mocker.Header{
		"Authorization": tester.Token,
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user.Profile{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	assert.Equal(t, profile.Id, tester.Uid)
	assert.Equal(t, *profile.Email, tester.Username)
}
