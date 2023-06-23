package ratelimiter_test

import (
	"testing"
	"time"
	"track-driver-performance/ratelimiter"
)

func TestSlidingWindowRatelimiter(t *testing.T) {
	limiter := ratelimiter.NewRateLimiter("sliding-window", 5 , 5, time.Second)

	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request should have been allowed")
		}

	}
	if limiter.Allow() {
		t.Errorf("Request should be blocked ")
	}
	time.Sleep(time.Second)

	if !limiter.Allow(){
		t.Errorf("request should be allowed")
	}

}


func TestTokenBucketRatelimiter(t *testing.T) {
	limiter := ratelimiter.NewRateLimiter("token-bucket", 5 , 5, time.Second)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request should have been allowed")
		}

	}
	if limiter.Allow() {
		t.Errorf("Request should be blocked ")
	}
	time.Sleep(time.Second)

	if !limiter.Allow(){
		t.Errorf("request should be allowed")
	}
}





