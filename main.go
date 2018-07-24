package main

import (
	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
	"log"
	"strings"
)

func main() {
	args := getArgs()
	options := redis.Options{
		Addr: strings.Join([]string{args.redisHost, args.redisPort}, ":"),
		DB:   0,
	}

	client := redis.NewClient(&options)
	tokenContext := NewTokenContext(client, args.redisPrefix)

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(client, args.rateLimit))

	router.POST("/graphql", createGraphQLHandler(tokenContext, args))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.httpPort}, "")))
}
