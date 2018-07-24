package main

import (
	"encoding/json"
	"errors"
	redis "github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"strings"
)

// TokenGrant describes the properties that a token represents
type TokenGrant struct {
	ContainerID string `json:"containerID"`
	ImageTag    string `json:"imageTag"`
}

// TokenContext describes the context needed for managing tokens
type TokenContext struct {
	redisClient *redis.Client
	keyPrefix   string
}

// NewTokenContext creates a new context that enables managing tokens
func NewTokenContext(client *redis.Client, prefix string) TokenContext {
	return TokenContext{
		redisClient: client,
		keyPrefix:   prefix,
	}
}

func createTokenKey(prefix string, token string) string {
	return strings.Join([]string{prefix, token}, "_")
}

func (c TokenContext) GetGrant(token string) (*TokenGrant, error) {
	key := createTokenKey(c.keyPrefix, token)
	val, err := c.redisClient.Get(key).Result()

	if err != nil {
		return nil, errors.New("invalid access token")
	}

	var grant TokenGrant
	err = json.Unmarshal([]byte(val), &grant)

	if err != nil {
		return nil, errors.New("invalid token grant")
	}

	return &grant, nil
}

func (c TokenContext) CreateToken(grant *TokenGrant) (string, error) {
	token := uuid.NewV4().String()
	marshalledGrant, err := json.Marshal(grant)
	if err != nil {
		return "", err
	}

	c.redisClient.Set(
		createTokenKey(c.keyPrefix, token),
		string(marshalledGrant),
		0,
	)

	return token, nil
}

func (c TokenContext) RevokeToken(token string) {
	c.redisClient.Del(
		createTokenKey(c.keyPrefix, token),
	)
}
