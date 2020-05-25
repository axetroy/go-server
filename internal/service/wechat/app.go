package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type Response struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

func FetchOpenID(code string) (*Response, error) {
	// 如果是测试环境，则返回一个写死的数据，方便测试
	if dotenv.Test {
		return &Response{
			OpenID: "oPl3r0AZdJxd7fO0HhMb99Te1311",
		}, nil
	}

	c := model.Config{Name: model.ConfigFieldNameWechatApp.Field}

	if err := database.Db.Model(&c).Where(&c).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoConfig
		}
		return nil, err
	}

	wechatConfig := model.ConfigFieldWechatApp{}

	if err := json.Unmarshal([]byte(c.Fields), &wechatConfig); err != nil {
		return nil, err
	}

	wechatUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", wechatConfig.AppID, wechatConfig, code)

	r, reqErr := http.Get(wechatUrl)

	if reqErr != nil {
		return nil, reqErr
	}

	resBytes, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		return nil, ioErr
	}

	reqRes := Response{}

	if jsonErr := json.Unmarshal(resBytes, &reqRes); jsonErr != nil {
		return nil, jsonErr
	}

	if reqRes.ErrMsg != "" {
		return nil, exception.Unknown.New(reqRes.ErrMsg)
	}

	return &reqRes, nil
}
