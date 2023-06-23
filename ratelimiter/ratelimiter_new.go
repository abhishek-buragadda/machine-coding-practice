package ratelimiter

import "time"

/*
	2:40 PM

Design a ratelimiter.


Ratelimiter is a layer before the service, which  can be used to limit the no of calls happening to the service to protect
it from DDOS type of attack or limit the bandwidth for a particular user.


Different strategies for ratelimiteer

	- Sliding window rate limiter .

			- give a window interval, we will fix the number of requests coming in that interval.

	- Token bucket ratelimiter .
			- we have a bucket full of tokens, i.e fixed amount of tokens and we let the request go through,
if we have a token avialable else block. We have a rate of refresh of the tokens.


*/


type Strategy interface {
	IsAllow() bool
}

type SlidingWindowStrategy struct {
	windowInterval time.Duration
	lastResetTime time.Time
	tokenCount int
	currTokenCount int

}
func(s *SlidingWindowStrategy) IsAllow() bool {
	if s.lastResetTime.Add(s.windowInterval).After(time.Now()){
			if s.currTokenCount > 0 {
				s.currTokenCount--
				return true
			}
			return false
	}
	s.lastResetTime = time.Now()
	s.currTokenCount = s.tokenCount-1
	return true
}

type TokenBucketStrategy struct {
	tokenRateInSec int
	currentTokens int
	totalTokens int
	lastResetTime time.Time


}

func(t *TokenBucketStrategy) IsAllow() bool {

	timeElapsed := time.Now().Sub(t.lastResetTime)
	tokenToAdd := int(timeElapsed.Seconds()) * t.tokenRateInSec

	if tokenToAdd > 0 {
		t.currentTokens =  t.currentTokens+tokenToAdd
		if t.currentTokens > t.totalTokens {
			t.currentTokens = t.totalTokens
		}
		t.lastResetTime = time.Now()
	}

	if t.currentTokens > 0 {
		t.currentTokens--
		return true

	}
	return false
}



type RatelimiterNew struct {
	strategy Strategy

}

func(r *RatelimiterNew) Allow() bool {
	return r.strategy.IsAllow()
}


func main(){



}