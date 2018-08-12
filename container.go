package main

import (
	"github.com/fsouza/go-dockerclient"
	"strings"
)

func environmentContainsHostname(variables []string, hostnames []string) bool {
	for _, variable := range variables {
		res := strings.SplitN(variable, "=", 2)

		if len(res) < 2 {
			continue
		} else if res[0] == "VIRTUAL_HOST" || res[0] == "LETSENCRYPT_HOST" {
			for _, hostname := range hostnames {
				for _, assignedHostname := range strings.Split(res[1], ",") {
					if hostname == assignedHostname {
						return true
					}
				}
			}
		}
	}

	return false
}

func hostnameAssigned(client *docker.Client, hostname string) (assignedContainer *docker.Container, err error) {
	hostnames := strings.Split(hostname, ",")
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	for i := range containers {
		container, err := client.InspectContainer(containers[i].ID)
		if err != nil {
			return nil, err
		}

		if environmentContainsHostname(container.Config.Env, hostnames) {
			return container, nil
		}
	}

	return nil, nil
}
