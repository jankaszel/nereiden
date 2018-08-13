package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func requireToken(token string) gin.HandlerFunc {
	requiredHeader := "Token " + token

	return func(c *gin.Context) {
		if queryToken, exists := c.GetQuery("access_token"); exists && queryToken == token {
			return
		} else if c.GetHeader("Authorization") == requiredHeader {
			return
		}

		c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}
