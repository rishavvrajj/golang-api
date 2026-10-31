package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "9001", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "Hello from backend on port %s\n", *port)
	})

	log.Printf("Backend listening on :%s", *port)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}