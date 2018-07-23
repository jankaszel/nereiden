package main

import (
	"encoding/json"
	"errors"
	redis "github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"strings"
)

const prefix string = "token"

// TokenConf describes the properties that a token represents
type TokenConf struct {
	ContainerID string `json:"containerID"`
	ImageTag    string `json:"imageTag"`
}

func createTokenKey(token string) string {
	return strings.Join([]string{prefix, token}, "_")
}

func acquireConf(client *redis.Client, token string) (tokenConf *TokenConf, err error) {
	key := createTokenKey(token)
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

func createToken(client *redis.Client, conf *TokenConf) (string, error) {
	token := uuid.NewV4().String()
	marshalledConf, err := json.Marshal(conf)
	if err != nil {
		return "", err
	}

	client.Set(createTokenKey(token), string(marshalledConf), 0)

	return token, nil
}

func revokeToken(client *redis.Client, token string) {
	key := createTokenKey(token)
	client.Del(key)
}
