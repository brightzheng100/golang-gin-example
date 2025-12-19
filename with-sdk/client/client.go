package main

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	instana "github.com/instana/go-sensor"
	"github.com/instana/go-sensor/instrumentation/instalogrus"
	"github.com/opentracing/opentracing-go/ext"
)

var server_url string
var collector instana.TracerLogger

func init() {
	////////////////////////////////////////////////////
	// create an instana collector
	collector = instana.InitCollector(&instana.Options{
		Service: "http-client",
		Tracer:  instana.DefaultTracerOptions(),
	})

	// Register instalogrus hook within the global logger
	log.AddHook(instalogrus.NewHook(collector))
	////////////////////////////////////////////////////
}

func main() {
	flag.StringVar(&server_url, "s", "http://localhost:8080", "Server address")
	flag.Parse()

	if server_url == "" {
		flag.Usage()
		os.Exit(2)
	}

	var endpoints = []string{"/bing", "/v1/bing", "/v1/error"}

	client := &http.Client{
		Transport: instana.RoundTripper(collector, nil),
	}

	for {
		// Make the GET request
		endpoint := endpoints[rand.Intn(len(endpoints))]
		log.Infof("endpoint: %s", endpoint)
		url := server_url + endpoint

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatalf("failed to create request: %s", err)
		}

		span := collector.Tracer().StartSpan("http-client-request")
		span.SetTag(string(ext.SpanKind), "entry")
		span.SetTag("endpoint", endpoint)
		ctx := instana.ContextWithSpan(context.Background(), span)
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			//span.LogKV("endpoint", endpoint, "error", err)
			log.WithContext(ctx).Fatalf("failed to GET %v: %v", endpoint, err)
		}
		span.Finish()

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
