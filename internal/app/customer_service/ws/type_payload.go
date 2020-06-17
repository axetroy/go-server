// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package ws

type MessageTextPayload struct {
	Text string `json:"text" validate:"required,max=255" comment:"消息体"`
}

type MessageImagePayload struct {
	Image string `json:"image" validate:"required,url,max=255" comment:"图片URL"`
}

type QueuePayload struct {
	Location uint `json:"location" validate:"required,int,min=0" comment:"位置"`
}

type AuthPayload struct {
	Token string `json:"token" validate:"required,min=0" comment:"身份令牌"`
}

type RatePayload struct {
	Rate uint `json:"rate" validate:"required,int,min=1,max=5" comment:"评分"`
}
