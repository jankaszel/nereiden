package main

import (
	"errors"
	"github.com/graphql-go/graphql"
)

func createTokenRevocationMutation(tokenContext TokenContext) *graphql.Field {
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

			tokenContext.RevokeToken(token)
			return true, nil
		},
	}
}
