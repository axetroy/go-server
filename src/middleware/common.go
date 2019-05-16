package middleware

import (
	"github.com/gin-gonic/gin"
)

func Common(context *gin.Context) {
	header := context.Writer.Header()
	// alone dns prefect
	header.Set("X-DNS-Prefetch-Control", "on")
	// IE No Open
	header.Set("X-Download-Options", "noopen")
	// not cache
	header.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	header.Set("Expires", "max-age=0")
	// Content Security Policy
	header.Set("Content-Security-Policy", "default-src 'self'")
	// xss protect
	// it will caught some problems is old IE
	header.Set("X-XSS-Protection", "1; mode=block")
	// Referrer Policy
	header.Set("Referrer-Header", "no-referrer")
	// cros frame, allow same origin
	header.Set("X-Frame-Options", "SAMEORIGIN")
	// HSTS
	header.Set("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
	// no sniff
	header.Set("X-Content-Type-Options", "nosniff")
}
