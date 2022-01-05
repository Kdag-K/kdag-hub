package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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

// RemoveContainer removes a container.
func RemoveContainer(cli *client.Client, containerID string, force, removelinks, removevolumes bool) error {
	ctx := context.Background()
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force:         force,
		RemoveLinks:   removelinks,
		RemoveVolumes: removevolumes,
	}); err != nil {
		return err
	}
	
	return nil
}

// GetContainers lists containers.
func GetContainers(cli *client.Client, output bool) (map[string]string, error) {
	
	rtn := make(map[string]string)
	
	ctx := context.Background()
	arrRes, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		fmt.Println(err.Error())
		return rtn, err
	}
	
	for _, net := range arrRes {
		if output {
			fmt.Printf("%s   %s  %s\n", net.Names[0], net.ID, net.Status)
		}
		if len(net.Names) > 0 {
			rtn[strings.TrimLeft(net.Names[0], "/")] = net.ID
		}
	}
	
	return rtn, nil
}

// StopContainer stops a container
func StopContainer(cli *client.Client, containerID string) error {
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		return err
	}
	return nil
}

// ConnectContainerToNetwork connects a created container to an extant network.
func ConnectContainerToNetwork(cli *client.Client, networkID string, containerID string, ip string) error {
	
	ctx := context.Background()
	
	return cli.NetworkConnect(ctx, networkID, containerID,
		&network.EndpointSettings{
			IPAMConfig: &network.EndpointIPAMConfig{IPv4Address: ip},
			IPAddress:  ip,
		})
}

// CopyToContainer copies a directory / file to a container
func CopyToContainer(cli *client.Client, containerID, localSrcPath, containerDestPath string) error {
	ctx := context.Background()
	
	archive, err := newTarArchiveFromPath(localSrcPath)
	if err != nil {
		return err
	}
	
	err = cli.CopyToContainer(ctx, containerID, containerDestPath, archive, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}
	
	return nil
}

// CopyFromContainer copies a file / directory from a container.
func CopyFromContainer(cli *client.Client, containerID, containerSrcPath, localDestPath string) error {
	
	ctx := context.Background()
	
	content, stat, err := cli.CopyFromContainer(ctx, containerID, containerSrcPath)
	if err != nil {
		return err
	}
	defer content.Close()
	
	srcInfo := archive.CopyInfo{
		Path:       containerSrcPath,
		Exists:     true,
		IsDir:      stat.Mode.IsDir(),
		RebaseName: "",
	}
	
	preArchive := content
	return archive.CopyTo(preArchive, srcInfo, localDestPath)
}

func newTarArchiveFromPath(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	
	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		
		if fi.IsDir() {
			return nil
		}
		
		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}
		
		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})
	
	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}