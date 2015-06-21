package main

import (
	"log"
    "github.com/samalba/dockerclient"
)

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func getContainersWithAddresses(targetedContainerNames []string) map[string]string  {
    // Init the client
    docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

    // Get only running containers
    containers, err := docker.ListContainers(false, false, "")
    if err != nil {
        log.Fatal(err)
    }

    m := make(map[string]string)

    for _,container := range containers {
    	id := container.Id
    	info, _ := docker.InspectContainer(id)
    	if stringInSlice(info.Name[1:], targetedContainerNames) {
    		m[info.Name] = info.NetworkSettings.IPAddress
    	}
    }

    return m
}