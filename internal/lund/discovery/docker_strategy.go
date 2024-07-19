package discovery

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/podikoglou/lund/internal/lund"
)

type DockerDiscoveryStrategy struct {
	client  *client.Client
	filters filters.Args
}

func NewDockerDiscoveryStrategy() DockerDiscoveryStrategy {
	// create client using settings from env variables
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	// create filters for searching for the containers with a certain label
	filters := filters.NewArgs()
	filters.Add("label", "lund.enable=true")

	return DockerDiscoveryStrategy{
		client:  client,
		filters: filters,
	}
}

func (s DockerDiscoveryStrategy) Discover() []*lund.Server {
	containers, err := s.client.ContainerList(
		context.TODO(),

		container.ListOptions{
			Filters: s.filters,
		},
	)

	if err != nil {
		log.Fatal(err) // should this be fatal?
	}

	var servers []*lund.Server

	for _, container := range containers {
		// try to figure out the port
		var port uint16

		val, exists := container.Labels["lund.port"]

		if exists {
			parsed, err := strconv.Atoi(val)

			if err != nil {
				log.Fatalf("invalid port: %s", val)
			} else {
				port = uint16(parsed)
			}
		} else {
			// if the lund.port label doesn't exist, try to guess it based on the first open port
			port = container.Ports[0].PublicPort
		}

		// try to hopefully figure out the host (it *might* not be this one)
		host := container.NetworkSettings.Networks["lund"].IPAddress

		// craft URL
		// TODO: add another label for a path prefix, or host to be used in the
		// http header.
		url := fmt.Sprintf("http://%s:%s", host, port)

		// add the server to the list
		servers = append(servers, &lund.Server{
			URL:   url,
			Alive: atomic.Bool{},
		})

	}

	return servers
}
