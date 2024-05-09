package ratelimit

import (
	"github.com/juju/ratelimit"
)

type RateLimiterInterceptor struct {
	TokenBucket *ratelimit.Bucket
}

func New(rate float64, capacity int64) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		TokenBucket: ratelimit.NewBucketWithRate(rate, capacity),
	}
}

func (r *RateLimiterInterceptor) Limit() bool {
	tokenRes := r.TokenBucket.TakeAvailable(1)
	return tokenRes == 0
}
