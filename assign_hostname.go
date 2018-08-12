package main

import (
	"errors"
	"github.com/falafeljan/docker-recreate"
	"github.com/fsouza/go-dockerclient"
	"github.com/graphql-go/graphql"
)

var recreationResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContainerRecreationResponse",
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

func assignContainerHostname(
	containerID string,
	hostname string,
	letsEncryptEmail string,
) (*recreate.Recreation, error) {
	options := recreate.DockerOptions{
		PullImage:       false,
		DeleteContainer: true,
		Registries:      []recreate.RegistryConf{},
	}

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return nil, err
	}

	context := recreate.NewContextWithClient(options, client)

	env := make(map[string]string)
	env["VIRTUAL_HOST"] = hostname
	env["LETSENCRYPT_HOST"] = hostname
	env["LETSENCRYPT_EMAIL"] = letsEncryptEmail

	recreation, err := context.Recreate(
		containerID,
		"",
		recreate.ContainerOptions{
			Env: env,
		},
	)

	if err != nil {
		return nil, err
	}

	err = client.ConnectNetwork("nginx-proxy", docker.NetworkConnectionOptions{
		Container: recreation.NewContainerID,
	})

	if err != nil {
		return nil, err
	}

	return recreation, nil
}

func createHostnameAssignmentMutation(letsEncryptEmail string) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewNonNull(recreationResponseType),
		Args: graphql.FieldConfigArgument{
			"containerID": &graphql.ArgumentConfig{
				Description: "The ID of the container",
				Type:        graphql.NewNonNull(graphql.String),
			},
			"hostname": &graphql.ArgumentConfig{
				Description: "The hostname that is going to be assigned to the specified container.",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			containerID, ok := params.Args["containerID"].(string)
			if !ok {
				return nil, errors.New("`containerID` is expected to be a string")
			}

			hostname, ok := params.Args["hostname"].(string)
			if !ok {
				return nil, errors.New("`hostname` is expected to be a string")
			}

			return assignContainerHostname(containerID, hostname, letsEncryptEmail)
		},
	}
}
