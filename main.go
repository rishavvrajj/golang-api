package main

import (
	"log"
	"net/http"
)

func main() {
	addresses := []string{
		"https://example.com",
		"https://httpbin.org",
		"https://api.github.com",
	}

	servers := make([]Server, 0, len(addresses))

	for _, addr := range addresses {
		server, err := NewSimpleServer(addr)
		if err != nil {
			log.Fatalf("invalid backend address %q: %v", addr, err)
		}

		servers = append(servers, server)
	}

	lb := NewLoadBalancer("8000", servers)

	log.Printf("Load balancer listening on :%s", lb.Port())

	if err := http.ListenAndServe(":"+lb.Port(), lb); err != nil {
		log.Fatal(err)
	}
}