package main

import (
	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"

	limiter "github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

func limiterMiddleware(client *redis.Client) (middleware gin.HandlerFunc) {
	rate, err := limiter.NewRateFromFormatted("1-S")
	if err != nil {
		panic(err)
	}

	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter",
		MaxRetry: 3,
	})
	if err != nil {
		panic(err)
	}

	middleware = mgin.NewMiddleware(
		limiter.New(store, rate))

	return middleware
}
