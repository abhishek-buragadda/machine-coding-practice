package load_balancer

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

/*
		load balancer via code
		Design a load balancer which will be before bunch of servers and distribute the load among them.

	 	Ability to decide on the strategy for distributing the load.
		- Round Robin
		- Random
		- Weighted Round Robin.

*/

type Strategy interface {
	GetNextServer(servers []string) string
}

type RoundRobinStrategy struct {
	curr int
}

func NewStrategy(strategy string) Strategy {
	if strategy == "round-robin" {
		return &RoundRobinStrategy{}
	}
	return &RandomStrategy{}
}

func (r *RoundRobinStrategy) GetNextServer(servers []string) string {
	selectedServer := servers[r.curr]
	r.curr = (r.curr + 1) % len(servers)
	return selectedServer
}

type RandomStrategy struct{}

func (randomStrategy *RandomStrategy) GetNextServer(servers []string) string {
	index := rand.Intn(len(servers))
	return servers[index]
}

type LoadBalancer struct {
	lock     sync.Mutex
	servers  []string
	strategy Strategy
}

func NewLoadBalancer(servers []string, strategy Strategy) *LoadBalancer {
	return &LoadBalancer{
		lock:     sync.Mutex{},
		servers:  servers,
		strategy: strategy,
	}
}

func (l *LoadBalancer) AddServer(host string) {
	l.lock.Lock()
	l.servers = append(l.servers, host)
	l.lock.Unlock()
}

func (l *LoadBalancer) RemoveServer(host string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	index := -1
	for i := 0; i < len(l.servers); i++ {
		if l.servers[i] == host {
			index = i
		}
	}
	l.servers = append(l.servers[:index], l.servers[index+1:]...)

}

func(l *LoadBalancer) GetServerCount() int {
	return len(l.servers)
}

func (l *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
	serverHost := l.strategy.GetNextServer(l.servers)
	fmt.Printf("Forwarding the request to the server: %s", serverHost)
	http.Redirect(w, r, serverHost+r.URL.String(), http.StatusTemporaryRedirect)

}

func main() {
	l := LoadBalancer{
		servers:  []string{"one", "two"},
		strategy: NewStrategy("round-robin"),
	}
	http.HandleFunc("/", l.handleRequest)
	fmt.Println("Load balancer listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
