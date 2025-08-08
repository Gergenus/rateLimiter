package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	burst  int
	bucket chan Token
	limit  Limit
}

type Limit struct {
	speed int
}

type Token struct{}

func NewToken() Token {
	return Token{}
}

func NewRateLimiter(limit Limit, burst int) RateLimiter {
	bucket := make(chan Token, burst)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(speed int) {
		wg.Done()
		for {
			<-time.Tick(time.Duration(speed) * time.Second)
			bucket <- NewToken()
		}
	}(limit.speed)
	wg.Wait()
	return RateLimiter{
		limit:  limit,
		burst:  burst,
		bucket: bucket,
	}
}

func NewLimit(speed int) Limit {
	return Limit{
		speed: speed,
	}
}

func (r *RateLimiter) Allow() bool {
	select {
	case <-r.bucket:
		return true
	default:
		return false
	}
}
