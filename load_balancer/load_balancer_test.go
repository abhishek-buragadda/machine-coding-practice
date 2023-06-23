package load_balancer_test

import (
	"testing"
	"track-driver-performance/load_balancer"
)
import "github.com/stretchr/testify/assert"

func TestRoundRobinStrategy_GetNextServer(t *testing.T) {
	servers := []string{"one", "two", "three"}
	strategy := load_balancer.NewStrategy("round-robin")
	nextServer := strategy.GetNextServer(servers)
	assert.Equal(t, nextServer, "one")
	nextServer = strategy.GetNextServer(servers)
	assert.Equal(t, nextServer, "two")

	nextServer = strategy.GetNextServer(servers)
	assert.Equal(t, nextServer, "three")

	nextServer = strategy.GetNextServer(servers)
	assert.Equal(t, nextServer, "one")

}

func TestRandomStrategy_GetNextServer(t *testing.T) {
	servers := []string{"one", "two", "three"}
	strategy := load_balancer.NewStrategy("random")
	nextServer := strategy.GetNextServer(servers)
	isServerFound := false
	for _, server := range servers {
		if server == nextServer {
			isServerFound = true
		}

	}
	if !isServerFound {
		t.Errorf("Unable to find the server ")
	}

}

func TestLoadBalancerAddServer(t *testing.T) {

	l := load_balancer.NewLoadBalancer([]string{"one", "two"}, load_balancer.NewStrategy("round-robin"))

	l.AddServer("three")
	assert.Equal(t, l.GetServerCount(), 3)
}

func TestLoadBalancerRemoveServer(t *testing.T) {

	l := load_balancer.NewLoadBalancer([]string{"one", "two"}, load_balancer.NewStrategy("round-robin"))

	l.RemoveServer("two")
	assert.Equal(t, l.GetServerCount(), 1)
}
