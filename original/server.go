// (c) Copyright IBM Corp. 2021
// (c) Copyright Instana Inc. 2016

package main

import (
	"errors"
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	engine.GET("/ping", func(c *gin.Context) {
		log.WithFields(log.Fields{
			"endpoint": "/ping",
		}).Info("/ping endpoint called")

		c.JSON(200, gin.H{
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
			"message": "pong",
		})
	})

	v1.GET("/error", func(c *gin.Context) {
		log.WithFields(log.Fields{
			"endpoint": "/v1/error",
		}).Error("/v1/error endpoint called")

		c.Error(errors.New("something went wrong"))
	})

	engine.Run(listenAddr)
}
