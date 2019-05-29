package system_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller/system"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	r := system.GetSystemInfo()

	assert.Equal(t, 1, r.Status)
	assert.Equal(t, "", r.Message)
}

func TestGetSystemInfoRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/system", nil, &header)

	res := schema.Response{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	n := system.Info{}

	assert.Nil(t, tester.Decode(res.Data, &n))
}
