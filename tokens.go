package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
	"net/http"
	"strings"
)

const prefix string = "token"

// TokenConf describes the properties a token represents
type TokenConf struct {
	ContainerID string `json:"containerId"`
	ImageTag    string `json:"imageTag"`
}

func acquireConf(client *redis.Client, token string) (tokenConf *TokenConf, err error) {
	key := strings.Join([]string{prefix, token}, "_")
	val, err := client.Get(key).Result()

	if err != nil {
		return nil, errors.New("invalid access token")
	}

	var conf TokenConf
	err = json.Unmarshal([]byte(val), &conf)

	if err != nil {
		return nil, errors.New("invalid token configuration")
	}

	return &conf, nil
}

func tokenMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("access_token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "no `access_token` provided"})

			return
		}

		tokenConf, err := acquireConf(client, token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": err.Error()})

			return
		}

		c.Set("tokenConf", *tokenConf)
		c.Next()
	}
}
