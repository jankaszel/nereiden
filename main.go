package main

import (
	"fmt"
	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
	"log"
	"strings"
)

func main() {
	args := getArgs()
	options := redis.Options{
		Addr: strings.Join([]string{args.redisHost, args.redisPort}, ":"),
		DB:   0}
	fmt.Printf("%+v\n", args)
	client := redis.NewClient(&options)

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(client))

	router.GET(
		"/recreate",
		tokenMiddleware(client),
		createHandler(args.registries))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.httpPort}, "")))
}
