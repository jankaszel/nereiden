package main

import (
	"encoding/json"
	"errors"
	redis "github.com/go-redis/redis"
	"strings"
)

const prefix string = "token"

// TokenConf describes the properties that a token represents
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
