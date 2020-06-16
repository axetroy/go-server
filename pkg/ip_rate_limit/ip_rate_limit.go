package ip_rate_limit

import (
	"golang.org/x/time/rate"
	"sync"
)

// IPRateLimiter .
type IPRateLimiter struct {
	sync.RWMutex
	ips map[string]*rate.Limiter
	r   rate.Limit
	b   int
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.Lock()
	defer i.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.Unlock()
		return i.AddIP(ip)
	}

	i.Unlock()

	return limiter
}
