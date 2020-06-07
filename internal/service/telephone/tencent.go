// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package telephone

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func NewTencent() *Tencent {
	c := &Tencent{}

	return c
}

type Tencent struct {
}

type tencentTel struct {
	Mobile     string `json:"mobile"`      // 手机号
	NationCode string `json:"nation_code"` // 国家/地区码
}

// 文档相关: https://cloud.tencent.com/document/product/382/5976
type tencentCloudParams struct {
	Ext    *string    `json:"ext"`    // 用户的 session 内容，腾讯 server 回包中会原样返回，可选字段，不需要时设置为空
	Extend *string    `json:"extend"` // 短信码号扩展号，格式为纯数字串，其他格式无效。默认没有开通，如需开通请联系 sms helper
	Params []string   `json:"params"` // 模板参数，具体使用方法请参见下方说明。若模板没有参数，请设置为空数组
	Sig    string     `json:"sig"`    // App 凭证，具体计算方式请参见下方说明
	Sign   string     `json:"sign"`   // 短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名
	Tel    tencentTel `json:"tel"`    // 国际电话号码，格式依据 e.164 标准为：+[国家（或地区）码][手机号] ，示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号
	Time   int64      `json:"time"`   // 请求发起时间，UNIX 时间戳（单位：秒），如果和系统时间相差超过 10 分钟则会返回失败
	TplId  int        `json:"tpl_id"` // 模板 ID，必须填写已审核通过的模板 ID
}

type tencentCloudResponse struct {
	Result int     `json:"result"` // 错误码，0表示成功（计费依据），非0表示失败，更多详情请参见 错误码
	ErrMsg string  `json:"errmsg"` // 错误消息，result 非0时的具体错误信息
	Ext    *string `json:"ext"`    // 用户的 session 内容，腾讯 server 回包中会原样返回
	Fee    *int    `json:"fee"`    // 短信计费的条数，计费规则请参考 国内短信内容长度计算规则 或 国际/港澳台短信内容长度计算规则
	Sid    *string `json:"sid"`    // 本次发送标识 ID，标识一次短信下发记录
}

func (c *Tencent) getAuthTemplateID() string {
	return config.Telephone.Tencent.TemplateCodeAuth
}

func (c *Tencent) getResetPasswordTemplateID() string {
	return config.Telephone.Tencent.TemplateCodeResetPassword
}

func (c *Tencent) getRegisterTemplateID() string {
	return config.Telephone.Tencent.TemplateCodeRegister
}

func (c *Tencent) send(phone string, templateID string, templateMap map[string]string) error {
	tplId, err := strconv.Atoi(templateID)

	if err != nil {
		return err
	}

	appKey := config.Telephone.Tencent.AppKey
	unixTIme := time.Now().Unix()
	randomStr := util.RandomNumeric(16)

	// 计算签名
	h := sha256.New()

	_, err = h.Write([]byte(fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", appKey, randomStr, unixTIme, phone)))

	if err != nil {
		return err
	}

	sig := h.Sum(nil)

	params := tencentCloudParams{
		Params: []string{},
		Sig:    string(sig),
		Sign:   config.Telephone.Tencent.Sign,
		Tel: tencentTel{
			Mobile:     phone,
			NationCode: "86",
		},
		Time:  unixTIme,
		TplId: tplId,
	}

	b, err := json.Marshal(params)

	if err != nil {
		return err
	}

	body := bytes.NewReader(b)

	r, err := http.Post(fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s", appKey, randomStr), "application/json", body)

	if err != nil {
		return exception.SendMsgFail
	}

	defer func() {
		_ = r.Body.Close()
	}()

	if r.StatusCode != http.StatusOK {
		return exception.SendMsgFail
	}

	resBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	res := tencentCloudResponse{}

	if err := json.Unmarshal(resBytes, &res); err != nil {
		return err
	}

	// 非 0 表示失败
	if res.Result != 0 {
		return exception.SendMsgFail
	}

	return nil
}

func (c *Tencent) SendAuthCode(phone string, code string) error {
	return c.send(phone, c.getAuthTemplateID(), map[string]string{
		"code": code,
	})
}

func (c *Tencent) SendResetPasswordCode(phone string, code string) error {
	return c.send(phone, c.getResetPasswordTemplateID(), map[string]string{
		"code": code,
	})
}

func (c *Tencent) SendRegisterCode(phone string, code string) error {
	return c.send(phone, c.getRegisterTemplateID(), map[string]string{
		"code": code,
	})
}
