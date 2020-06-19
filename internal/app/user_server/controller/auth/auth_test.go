package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/app/user_server/controller/user"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/pkg/proto"
	"github.com/axetroy/go-server/tester"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQRCodeGenerateLoginLink(t *testing.T) {
	r := auth.QRCodeGenerateLoginLink(helper.Context{
		Ip:        "127.0.0.1",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4173.0 Safari/537.36",
	})

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	assert.IsType(t, "", r.Data)

	link := r.Data.(string)

	p, err := proto.Parse(link)

	assert.Nil(t, err)

	b, err := p.Data()

	assert.Nil(t, err)

	var data map[string]interface{}

	assert.Nil(t, json.Unmarshal(b, &data))

	assert.NotEmpty(t, data["session_id"])
	assert.IsType(t, "", data["session_id"])

	assert.NotEmpty(t, data["expired_at"])
	assert.IsType(t, "", data["expired_at"])

	id, err := uuid.Parse(data["session_id"].(string))

	defer func() {
		_ = redis.QRCodeLoginCode.Del(context.Background(), id.String()).Err()
	}()

	assert.Nil(t, err)
	assert.Equal(t, data["session_id"], id.String())

	val, err := redis.QRCodeLoginCode.Get(context.Background(), id.String()).Result()

	assert.Nil(t, err)

	var body auth.QRCodeEntry

	assert.Nil(t, json.Unmarshal([]byte(val), &body))

	assert.Equal(t, "Intel Mac OS X 10_15_5", body.OS)
	assert.Equal(t, "Chrome", body.Browser)
	assert.Equal(t, "85.0.4173.0", body.Version)
	assert.Equal(t, "127.0.0.1", body.Ip)
}

func TestQRCodeLoginCheck(t *testing.T) {
	var link string

	// 先创建
	{
		r := auth.QRCodeGenerateLoginLink(helper.Context{
			Ip:        "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4173.0 Safari/537.36",
		})

		assert.IsType(t, "", r.Data)

		link = r.Data.(string)
	}

	// 再检查
	{
		r := auth.QRCodeLoginCheck(helper.Context{
			Ip:        "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4173.0 Safari/537.36",
		}, auth.QRCodeCheckParams{
			Url: link,
		})

		assert.Equal(t, exception.NoData.Code(), r.Status)
		assert.Equal(t, exception.NoData.Error(), r.Message)
		assert.Nil(t, r.Data)
	}

	p, err := proto.Parse(link)

	assert.Nil(t, err)

	b, err := p.Data()

	assert.Nil(t, err)

	var data map[string]interface{}

	assert.Nil(t, json.Unmarshal(b, &data))

	userInfo, err := tester.CreateUser()

	assert.Nil(t, err)

	defer func() {
		tester.DeleteUserByUid(userInfo.Id)
	}()

	{
		// 扫码登录
		r := user.QRCodeAuthGrant(helper.Context{
			Uid:       userInfo.Id,
			Ip:        "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4173.0 Safari/537.36",
		}, user.QRCodeAuthParams{
			Url: link,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Nil(t, r.Data)

		val, err := redis.QRCodeLoginCode.Get(context.Background(), data["session_id"].(string)).Result()

		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf(`{"os":"Intel Mac OS X 10_15_5","browser":"Chrome","version":"85.0.4173.0","ip":"127.0.0.1","user_id":"%s"}`, userInfo.Id), val)
	}

	// 再 check 一次
	{
		r := auth.QRCodeLoginCheck(helper.Context{
			Ip:        "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4173.0 Safari/537.36",
		}, auth.QRCodeCheckParams{
			Url: link,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		var profile schema.ProfileWithToken

		assert.Nil(t, r.Decode(&profile))

		assert.Equal(t, userInfo.Id, profile.Id)
		assert.Equal(t, userInfo.Username, profile.Username)
		assert.Equal(t, userInfo.Nickname, profile.Nickname)
	}

}
