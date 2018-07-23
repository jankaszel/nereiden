package main

import (
	"errors"
	redis "github.com/go-redis/redis"
	"github.com/graphql-go/graphql"
)

func createTokenRevocationMutation(client *redis.Client) *graphql.Field {
	return &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			"token": &graphql.ArgumentConfig{
				Description: "The to-be-revoked token.",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			token, ok := params.Args["token"].(string)
			if !ok {
				return nil, errors.New("`token` is expected to be a string")
			}

			revokeToken(client, token)
			return true, nil
		},
	}
}
