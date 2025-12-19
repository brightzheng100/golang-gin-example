package main

import (
	"flag"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var server_url string

func main() {
	flag.StringVar(&server_url, "s", "http://localhost:8080", "Server address")
	flag.Parse()

	if server_url == "" {
		flag.Usage()
		os.Exit(2)
	}

	var endpoints = []string{"/bing", "/v1/bing", "/v1/error"}

	for {
		// Make the GET request
		endpoint := server_url + endpoints[rand.Intn(len(endpoints))]
		log.Infof("Endpoint: %s", endpoint)

		resp, err := http.Get(endpoint)
		if err != nil {
			log.Fatalf("Error making GET request: %v", err)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Print the response status and body
		log.Infof("Response Status: %s", resp.Status)
		log.Infof("Response Body: %s", string(body))
		resp.Body.Close()

		// Sleep for a short duration before the next request
		time.Sleep(500 * time.Millisecond)
	}
}
