package main

import (
	gin "github.com/gin-gonic/gin"
	graphql "github.com/graphql-go/graphql"
	handler "github.com/graphql-go/handler"
)

func createGraphQLHandler(tokenContext TokenContext, args *Args) gin.HandlerFunc {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"recreateContainer": createContainerRecreationMutation(tokenContext, args.registries),
			"createToken":       createTokenCreationMutation(tokenContext),
			"revokeToken":       createTokenRevocationMutation(tokenContext),
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
		Schema: &schema,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
