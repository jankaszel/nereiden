package main

import (
	"github.com/fsouza/go-dockerclient"
	"strings"
)

func environmentContainsHostname(variables []string, hostname string) bool {
	for _, variable := range variables {
		res := strings.SplitN(variable, "=", 2)

		if len(res) < 2 {
			continue
		} else if (res[0] == "VIRTUAL_HOST" || res[0] == "LETSENCRYPT_HOST") && res[1] == hostname {
			return true
		}
	}

	return false
}

func hostnameAssigned(client *docker.Client, hostname string) (assignedContainer *docker.Container, err error) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	for i := range containers {
		container, err := client.InspectContainer(containers[i].ID)
		if err != nil {
			return nil, err
		}

		if environmentContainsHostname(container.Config.Env, hostname) {
			return container, nil
		}
	}

	return nil, nil
}
