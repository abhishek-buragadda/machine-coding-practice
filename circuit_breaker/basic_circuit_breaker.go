package circuit_breaker

import "time"



type BasicCircuitBreaker struct {
	state string
	errCount int64
	retryInterval time.Duration
}


func(b *BasicCircuitBreaker)  Execute( func() error ){

}
