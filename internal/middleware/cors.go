package middleware

import (
	"github.com/kataras/iris/v12"
	"net/http"
	"strings"
)

var (
	allowHeaders = strings.Join([]string{
		"accept",
		"origin",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"Content-Length",
		"Accept-Encoding",
		"Cache-Control",
		"X-CSRF-Token",
		"X-Requested-With",
		SignatureHeader,    // 接受签名的 Header
		PayPasswordHeader,  // 接收交易密码的 Header
		"X-Wechat-Binding", // 激活微信帐号
	}, ",")
	allowMethods = strings.Join([]string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}, ",")
)

func CORS() iris.Handler {
	return func(c iris.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Allow-Methods", allowMethods)

		if c.Request().Method == http.MethodOptions {
			c.StatusCode(http.StatusNoContent)
			c.EndRequest()
			return
		} else {
			c.Next()
		}
	}
}
