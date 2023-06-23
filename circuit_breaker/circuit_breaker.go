package circuit_breaker

import (
	"errors"
	"io"
	"net/http"
	"time"
)

/**
	Implement circuit breaker .  11:45 AM


	have configurable inputs based on which we can break the circuit.
	- error > x ms
	- errors (what are qualified as errors.
	- time_interval
	- retry_interval.
	- errorCode to return in case of matchbreak.
	- Should be a wrapper to the httpClient, fwd the request in case the condition is met else match break before.

	Config
		- error_interval
		- percentage_errors
		- errorCode

	Logic:
		- For each request coming to the circuitbreaker
			- count the request
			- categorize it to success/ failure.
			- calculate the percentage of sucecss/failure, if failure percentage> percentage configured matchbreak for xms.
			- reset the config after xms.



 */

type Config struct {
	url string
	errorPercentage float64
	errorInterval time.Duration
	errorCode  int
	timeout int
}

type RequestMetric struct {
	url string
	success int64
	failure int64
	total   int64
	isCircuitBroken bool

}

func(r *RequestMetric) Reset(){
	r.total =0
	r.failure = 0
	r.success = 0
}

type CircuitBreaker struct {
	errorStartTimeMap map[string]time.Time
	client http.Client
	requestMetricMap map[string]*RequestMetric
	config map[string]Config

}

func(c *CircuitBreaker) Do(req *http.Request)(*http.Response, error) {
	 _, ok :=c.requestMetricMap[req.URL.String()]
	 if !ok {
		c.requestMetricMap[req.URL.String()] = &RequestMetric{
			url:             req.URL.String(),
			success:         0,
			failure:         0,
			total:           1,
			isCircuitBroken: false,
		}
	}
	metric := c.requestMetricMap[req.URL.String()]
	if c.errorStartTimeMap[req.URL.String()].After(time.Now().Add(-1 * c.config[req.URL.String()].errorInterval)) {
		metric.Reset()
	}

	if metric.isCircuitBroken || c.isAboveThreshold(metric, req.URL.String()){
		return nil , errors.New("circuit broken")
	}
	resp, err :=  c.client.Do(req)
	if err!= nil {
		metric.failure = metric.failure+1
	}else{
		metric.Reset()
	}
	return resp, err
}

func(c *CircuitBreaker) isAboveThreshold(metric *RequestMetric, url string ) bool {
	errorPercentage := float64(metric.failure)/float64(metric.total)
	if errorPercentage> c.config[url].errorPercentage {
		c.errorStartTimeMap[url] = time.Now()
		metric.isCircuitBroken = true
		return true
	}
	return false

}

func (c *CircuitBreaker) Post(url, contentType string, body io.Reader) (resp *http.Response, err error){
	return c.client.Post(url, contentType, body)

}
func (c *CircuitBreaker) Head(url string) (resp *http.Response, err error){
	return c.client.Head(url)

}
func (c *CircuitBreaker) Get(url string) (resp *http.Response, err error){
	return c.client.Get(url)
}