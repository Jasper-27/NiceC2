package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// Load server key pair
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	// Create TLS config
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// InsecureSkipVerify: true,
	}

	// Create HTTPS server
	server := &http.Server{
		Addr:      ":8081",
		Handler:   http.HandlerFunc(nodeCheckIn),
		TLSConfig: config,
	}

	// Listen and serve HTTPS requests
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// Function for that runs each time the node checks in
func nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	var blank_response = string(`
    {
        "taskID": "0", 
        "task": "",  
        "arg" : ""
        
    }`)

	fmt.Fprint(w, blank_response)
}
