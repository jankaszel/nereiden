package main

import (
	gin "github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	args := getArgs()

	if args.InProduction {
		gin.SetMode("release")
	}

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(args.RateLimit))

	router.POST("/graphql", createGraphQLHandler(args.LetsEncryptEmail))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.HTTPPort}, "")))
}
