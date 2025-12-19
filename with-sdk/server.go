// (c) Copyright IBM Corp. 2021
// (c) Copyright Instana Inc. 2016

//go:build go1.17
// +build go1.17

package main

import (
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	instana "github.com/instana/go-sensor"
	"github.com/instana/go-sensor/instrumentation/instagin"
	"github.com/instana/go-sensor/instrumentation/instalogrus"
)

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "l", ":8080", "Server listen address")
	flag.Parse()

	if listenAddr == "" {
		flag.Usage()
		os.Exit(2)
	}

	engine := gin.Default()

	////////////////////////////////////////////////////
	// create an instana collector
	collector := instana.InitCollector(&instana.Options{
		Service: "gin-http-server",
		Tracer:  instana.DefaultTracerOptions(),
	})

	// Register instalogrus hook within the global logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.AddHook(instalogrus.NewHook(collector))

	// add middleware to the gin handlers
	instagin.AddMiddleware(collector, engine)
	////////////////////////////////////////////////////

	engine.GET("/ping", func(c *gin.Context) {
		log.WithFields(log.Fields{
			"endpoint": "/ping",
		}).Info("/ping endpoint called")

		c.JSON(200, gin.H{
			"success": true,
			"message": "pong",
		})
	})

	// use group: v1
	v1 := engine.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		log.WithFields(log.Fields{
			"endpoint": "/v1/ping",
		}).Info("/v1/ping endpoint called")

		c.JSON(200, gin.H{
			"success": true,
			"message": "pong",
		})
	})

	v1.GET("/error", func(c *gin.Context) {
		// simulate an error
		err := errors.New("something went wrong")

		log.WithContext(c.Request.Context()).WithFields(log.Fields{
			"endpoint": "/v1/error",
		}).Error(err)

		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
	})

	engine.Run(listenAddr)
}
