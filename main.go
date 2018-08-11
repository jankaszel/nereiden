package main

import (
	gin "github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	args := getArgs()

	if args.inProduction {
		gin.SetMode("release")
	}

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(args.rateLimit))

	router.POST("/graphql", createGraphQLHandler(args.letsEncryptEmail))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.httpPort}, "")))
}
