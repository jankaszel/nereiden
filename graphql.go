package main

import (
	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
	graphql "github.com/graphql-go/graphql"
	handler "github.com/graphql-go/handler"
)

func createGraphQLHandler(args *Args, client *redis.Client) gin.HandlerFunc {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"recreateContainer": createContainerRecreationMutation(args, client),
			"createToken":       createTokenCreationMutation(client),
			"revokeToken":       createTokenRevocationMutation(client),
		},
	})

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Mutation: rootMutation,
		Query:    rootQuery,
	})

	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
