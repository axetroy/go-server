// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package router

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net"
	"net/http"
)

type Context struct {
	context iris.Context
	err     error
}

func (c *Context) Param(key string) string {
	return c.context.Params().Get(key)
}

func (c *Context) Header(key string, value string) {
	c.context.Header(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.context.GetHeader(key)
}

func (c *Context) ClientIP() string {
	header := c.Request().Header

	// TODO: 客户端有可能设置这个头部用以欺骗服务端
	// 解决办法：添加白名单到 remoteAddr 中
	// 只认可白名单中的服务器发来的 `X-Real-Ip` 和 `X-Forwarded-For` 头部
	var clientIp = header.Get("X-Real-Ip")

	if clientIp == "" {
		clientIp = header.Get("X-Forwarded-For")
	}

	if clientIp == "" {
		addr := c.Request().RemoteAddr

		if ipStr, _, err := net.SplitHostPort(addr); err == nil {
			ip := net.ParseIP(ipStr)

			clientIp = string(ip)
		}
	}

	if clientIp == "" {
		clientIp = c.context.RemoteAddr()
	}

	return clientIp
}

func (c *Context) StatusCode(code int) {
	c.context.StatusCode(code)
}

func (c *Context) Request() *http.Request {
	return c.context.Request()
}

func (c *Context) ResetRequest(r *http.Request) {
	c.context.ResetRequest(r)
}

func (c *Context) Writer() http.ResponseWriter {
	return c.context.ResponseWriter()
}

func (c *Context) Application() context.Application {
	return c.context.Application()
}

func (c *Context) GetStatusCode() int {
	return c.context.GetStatusCode()
}

func (c *Context) GetBody() ([]byte, error) {
	return c.context.GetBody()
}

func (c *Context) ShouldBindJSON(pr interface{}) error {
	if err := c.context.ReadJSON(pr); err != nil {
		return exception.InvalidParams
	}
	return nil
}

func (c *Context) ShouldBindQuery(pr interface{}) error {
	if err := c.context.ReadQuery(pr); err != nil {
		err = exception.InvalidParams
	}
	return nil
}

func (c *Context) JSON(err error, data interface{}, meta *schema.Meta) {
	res := schema.Response{}

	if err != nil {
		res.Message = err.Error()

		if t, ok := err.(exception.Error); ok {
			res.Status = t.Code()
		} else {
			res.Status = exception.Unknown.Code()
		}
		res.Data = nil
		res.Meta = nil
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
		res.Meta = meta
	}

	c.Response(nil, res)
}

func (c *Context) ResponseFunc(err error, fn func() schema.Response) {
	if err != nil {
		res := schema.Response{}

		res.Message = err.Error()

		if t, ok := err.(exception.Error); ok {
			res.Status = t.Code()
		} else {
			res.Status = exception.Unknown.Code()
		}
		res.Data = nil
		res.Meta = nil

		_, _ = c.context.JSON(res)
	} else {
		_, _ = c.context.JSON(fn())
	}
}

func (c *Context) Response(err error, res schema.Response) {
	if err != nil {
		res.Message = err.Error()

		if t, ok := err.(exception.Error); ok {
			res.Status = t.Code()
		} else {
			res.Status = exception.Unknown.Code()
		}
		res.Data = nil
		res.Meta = nil
	}

	_, _ = c.context.JSON(res)
}

func (c *Context) Redirect(status int, url string) {
	c.context.Redirect(url, status)
}

func (c *Context) SetContext(key string, value interface{}) {
	c.context.Values().Set(key, value)
}

func (c *Context) GetContext(key string) interface{} {
	return c.context.Values().Get(key)
}

func (c *Context) Uid() string {
	return c.context.Values().GetString("uid")
}

func (c *Context) Next() {
	c.context.Next()
}

func Handler(handler func(c Context)) iris.Handler {
	return func(c iris.Context) {
		handler(Context{context: c})
	}
}
