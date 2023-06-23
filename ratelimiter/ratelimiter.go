package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type Ratelimiter struct {
	strategy string
	tokenRate int64
	interval time.Duration
	currTokens int64
	tokens int64
	lock sync.Mutex
	lastTokenResetTime time.Time
	rateAsTokenPerSec int64

}


func NewRateLimiter(strategy string, rateAsTokensPerSec, tokens  int64, interval time.Duration) *Ratelimiter {
	return &Ratelimiter{
		strategy:           strategy,
		tokenRate:          tokens,
		interval:           interval,
		currTokens:         tokens,
		tokens:             tokens,
		lock:               sync.Mutex{},
		lastTokenResetTime: time.Now(),
		rateAsTokenPerSec:  rateAsTokensPerSec,
	}

}
func(r *Ratelimiter) Allow() bool {
		r.lock.Lock()
		defer r.lock.Unlock()
		if r.strategy ==  "sliding-window" {
			return r.slidingWindow()
		}else if r.strategy == "token-bucket"{
			return r.tokenBucket()
		}
		return false

}

func(r *Ratelimiter) slidingWindow() bool {

	now := time.Now()
	//if now.Sub(r.lastTokenResetTime) >= r.interval {
	 if now.Add(-1 * r.interval).After(r.lastTokenResetTime) {
		r.currTokens =   r.tokens
		r.lastTokenResetTime = now
	}
	if r.currTokens > 0 {
		r.currTokens--
		return true
	}
	return false
}


func(r *Ratelimiter)  tokenBucket() bool{
	now := time.Now()
	elapsed := now.Sub(r.lastTokenResetTime)
	tokenToAdd := elapsed.Seconds() * float64(r.rateAsTokenPerSec)

	if tokenToAdd > 0 {
		r.currTokens = r.currTokens+ int64( tokenToAdd)

		if r.currTokens > r.tokens {
			r.currTokens = r.tokens
		}
		r.lastTokenResetTime = now
	}

	if r.currTokens > 0 {
		r.currTokens--
		return true
	}
	return false


}

func main() {
	ratelimiter := NewRateLimiter("sliding-window", 10, 10, time.Second)
	for i := 0; i < 15; i++ {
		if ratelimiter.Allow() {
			fmt.Println("ratelimiter allowed ")
		}else{
			fmt.Println("ratelimiter blocked")
		}
		time.Sleep(10* time.Millisecond)
	}
}

