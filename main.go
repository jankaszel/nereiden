package main

import (
	"log"

	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
)

func main() {
	option, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(option)

	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(client))

	router.GET("/recreate", tokenMiddleware(client), perform)
	log.Fatal(router.Run(":8090"))
}
