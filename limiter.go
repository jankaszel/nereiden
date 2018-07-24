package main

import (
	gin "github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	smemory "github.com/ulule/limiter/drivers/store/memory"
)

func limiterMiddleware(rateLimit string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		panic(err)
	}

	store := smemory.NewStore()
	middleware := mgin.NewMiddleware(
		limiter.New(store, rate))

	return middleware
}
