package docker

import (
	"context"
	"fmt"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetNetworks lists networks.
func GetNetworks(cli *client.Client, output bool) (map[string]string, error) {
	
	rtn := make(map[string]string)
	
	ctx := context.Background()
	arrRes, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return rtn, err
	}
	
	for _, net := range arrRes {
		if output {
			fmt.Printf("%s   %s  %s\n", net.Name, net.ID, net.Driver)
		}
		rtn[net.Name] = net.ID
	}
	
	return rtn, nil
}

// RemoveNetwork removes a network
func RemoveNetwork(cli *client.Client, networkID string) error {
	ctx := context.Background()
	return cli.NetworkRemove(ctx, networkID)
}