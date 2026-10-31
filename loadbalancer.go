package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	port            string
	mu              sync.Mutex
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:    port,
		servers: servers,
	}
}

func (lb *LoadBalancer) Port() string {
	return lb.port
}

func (lb *LoadBalancer) GetAvailableServer() (Server, bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return nil, false
	}

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.roundRobinCount%len(lb.servers)]
		lb.roundRobinCount++

		if server.IsAlive() {
			return server, true
		}

		log.Printf("Backend %s is unavailable", server.Address())
	}

	return nil, false
}

func (lb *LoadBalancer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	server, ok := lb.GetAvailableServer()
	if !ok {
		http.Error(
			rw,
			"No available backend servers",
			http.StatusServiceUnavailable,
		)
		return
	}

	fmt.Printf(
		"Forwarding %s %s to %s\n",
		req.Method,
		req.URL.Path,
		server.Address(),
	)

	server.Serve(rw, req)
}