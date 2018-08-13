package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func areAllOriginsAllowed(allowedOrigins []string) (bool, []string) {
	for i, origin := range allowedOrigins {
		if origin == "*" {
			return true, append(allowedOrigins[:i], allowedOrigins[i+1:]...)
		}
	}

	return false, allowedOrigins
}

func createCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allOriginsAllowed, remainingOrigins := areAllOriginsAllowed(allowedOrigins)

	return cors.New(cors.Config{
		AllowAllOrigins: allOriginsAllowed,
		AllowOrigins:    remainingOrigins,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	})
}
