package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
}

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func NewSimpleServer(addr string) (*SimpleServer, error) {
	serverURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse server URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(serverURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = serverURL.Host
	}

	return &SimpleServer{
		addr:  addr,
		proxy: proxy,
	}, nil
}

func (s *SimpleServer) Address() string {
	return s.addr
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (s *SimpleServer) IsAlive() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(s.addr)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < http.StatusInternalServerError
}