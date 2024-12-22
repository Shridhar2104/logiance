package ratelimit

import (
    "golang.org/x/time/rate"
    "time"
)

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), 1),
    }
}

func (r *RateLimiter) Allow() bool {
    return r.limiter.Allow()
}