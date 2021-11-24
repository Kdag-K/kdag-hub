package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/client"
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
