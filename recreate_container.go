package main

import (
	"errors"
	recreate "github.com/fallafeljan/docker-recreate/lib"
	redis "github.com/go-redis/redis"
	"github.com/graphql-go/graphql"
)

var recreationResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RecreationResponse",
	Fields: graphql.Fields{
		"previousContainerID": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the previous container",
		},
		"newContainerID": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the newly created container.",
		},
	},
})

func recreateContainer(args *Args, client *redis.Client, accessToken string) (*recreate.Recreation, error) {
	tokenConf, err := acquireConf(client, accessToken)

	if tokenConf.ContainerID == "" || tokenConf.ImageTag == "" {
		return nil, errors.New("internal error")
	}

	options := recreate.Options{
		PullImage:       true,
		DeleteContainer: true,
		Registries:      args.registries}

	recreation, err := recreate.Recreate(
		"unix:///var/run/docker.sock",
		tokenConf.ContainerID,
		tokenConf.ImageTag,
		&options)

	if err != nil {
		return nil, err
	}

	return recreation, nil
}

func createRecreateContainerMutation(args *Args, client *redis.Client) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewNonNull(recreationResponseType),
		Args: graphql.FieldConfigArgument{
			"accessToken": &graphql.ArgumentConfig{
				Description: "A token with respective grants",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			accessToken, ok := params.Args["accessToken"].(string)
			if !ok {
				return nil, errors.New("`accessToken` is expected to be a string")
			}

			return recreateContainer(args, client, accessToken)
		},
	}
}
