package docker

import (
	"context"
	
	"github.com/docker/docker/api/types"
	
	"github.com/docker/docker/client"
)

// StartContainer starts a container previously created by CreateContainerFromImage
func StartContainer(cli *client.Client, containerID string) error {
	
	ctx := context.Background()
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}