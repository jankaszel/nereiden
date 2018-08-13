package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func requireToken(token string) gin.HandlerFunc {
	requiredHeader := "Token " + token

	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != requiredHeader {
			c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
	}
}
