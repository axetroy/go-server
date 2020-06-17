// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/pkg/ip_rate_limit"
	"github.com/kataras/iris/v12"
	"net/http"
)

func RateLimit(maxConcurrencyPerMs uint) iris.Handler {
	limiter := ip_rate_limit.NewIPRateLimiter(5, int(maxConcurrencyPerMs))

	return func(c iris.Context) {
		limiter := limiter.GetLimiter(c.RemoteAddr())

		if !limiter.Allow() {
			c.StatusCode(http.StatusTooManyRequests)
			c.StopExecution()
			return
		}

		c.Next()
	}
}
