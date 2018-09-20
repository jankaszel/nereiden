package main

import (
	"errors"
	"fmt"
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
	proxyNetworkName string,
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

	container, err := client.InspectContainer(containerID)
	if err != nil {
		return nil, err
	}

	assignedContainer, err := hostnameAssigned(client, hostname)
	if err != nil {
		return nil, err
	} else if assignedContainer != nil && assignedContainer.ID != container.ID {
		return nil, fmt.Errorf(
			"Hostname `%s` is already assigned to container `%s`",
			hostname,
			assignedContainer.ID[:8],
		)
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

	newContainer, err := client.InspectContainer(recreation.NewContainerID)
	if err != nil {
		return nil, err
	}

	if _, ok := newContainer.NetworkSettings.Networks[proxyNetworkName]; !ok {
		err = client.ConnectNetwork(proxyNetworkName, docker.NetworkConnectionOptions{
			Container: recreation.NewContainerID,
		})
		if err != nil {
			return nil, err
		}
	}

	return recreation, nil
}

func createHostnameAssignmentMutation(
	proxyNetworkName string,
	letsEncryptEmail string,
) *graphql.Field {
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

			return assignContainerHostname(
				containerID,
				hostname,
				proxyNetworkName,
				letsEncryptEmail,
			)
		},
	}
}
