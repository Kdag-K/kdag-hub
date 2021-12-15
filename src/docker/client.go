package docker

import (
	"github.com/moby/moby/client"
)

// GetDockerClient returns a docker client.
func GetDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}