package docker

import (
	"context"
	"io"
	"os"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

// CreateContainerFromImage creates a container, returning its ID
func CreateContainerFromImage(cli *client.Client, imageName string, isImageRemote bool,
	nodeName string, cmd strslice.StrSlice, start bool) (string, error) {
	
	ctx := context.Background()
	
	if isImageRemote {
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		io.Copy(os.Stdout, out)
	}
	
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Cmd:   cmd,
		Image: imageName,
	}, nil, nil, nil, nodeName)
	if err != nil {
		return "", err
	}
	
	if start {
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return "", err
		}
	}
	//	fmt.Println(resp.ID)
	return resp.ID, nil
}

// StartContainer starts a container previously created by CreateContainerFromImage
func StartContainer(cli *client.Client, containerID string) error {
	
	ctx := context.Background()
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}