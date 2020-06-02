// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/kataras/iris/v12"
)

func Common() iris.Handler {
	return func(c iris.Context) {
		// alone dns prefect
		c.Header("X-DNS-Prefetch-Control", "on")
		// IE No Open
		c.Header("X-Download-Options", "noopen")
		// not cache
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Expires", "max-age=0")
		// Content Security Policy
		//c.Header("Content-Security-Policy", "default-src 'self'")
		// xss protect
		// it will caught some problems is old IE
		c.Header("X-XSS-Protection", "1; mode=block")
		// Referrer Policy
		c.Header("Referrer-Header", "no-referrer")
		// cros frame, allow same origin
		c.Header("X-Frame-Options", "SAMEORIGIN")
		// HSTS
		c.Header("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
		// no sniff
		c.Header("X-Content-Type-Options", "nosniff")

		c.Next()
	}
}
