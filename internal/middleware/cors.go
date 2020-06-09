package middleware

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"net/http"
)

var (
	allowHeaders = []string{
		"Accept",
		"Origin",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"Content-Length",
		"Accept-Encoding",
		"Cache-Control",
		"X-CSRF-Token",
		"X-Requested-With",
		SignatureHeader,   // 接受签名的 Header
		PayPasswordHeader, // 接收交易密码的 Header
	}
	allowMethods = []string{
		http.MethodOptions,
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
)

func CORS() iris.Handler {

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedMethods:   allowMethods,
		AllowedHeaders:   allowHeaders,
		MaxAge:           60,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	})

	return crs
}
