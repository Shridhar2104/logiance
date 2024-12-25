package ratelimit

import (
    "sync"
    "time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
    rate       float64
    maxTokens  float64
    tokens     float64
    lastUpdate time.Time
    mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified requests per second
func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{
        rate:       float64(rps),
        maxTokens:  float64(rps),
        tokens:     float64(rps),
        lastUpdate: time.Now(),
    }
}

// Allow checks if a request should be allowed based on the rate limit
func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(r.lastUpdate).Seconds()
    r.tokens = min(r.maxTokens, r.tokens+(elapsed*r.rate))
    r.lastUpdate = now

    if r.tokens >= 1 {
        r.tokens--
        return true
    }
    return false
}

// Wait blocks until a request is allowed based on the rate limit
func (r *RateLimiter) Wait() {
    for !r.Allow() {
        time.Sleep(time.Millisecond * 100)
    }
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// RateLimitedError represents a rate limit exceeded error
type RateLimitedError struct {
    Provider string
    Limit    int
}

func (e *RateLimitedError) Error() string {
    return "rate limit exceeded for provider " + e.Provider
}