package docker

import (
	"context"
	"fmt"
	
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
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

// CreateNetwork creates a docker network.
func CreateNetwork(cli *client.Client, networkName, subnet, iprange, gateway string) (string, error) {
	ctx := context.Background()
	
	common.DebugMessage("Creating options")
	opts := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		IPAM: &network.IPAM{
			Options: make(map[string]string),
			Config: []network.IPAMConfig{
				{
					Subnet:  subnet,
					IPRange: iprange,
					Gateway: gateway,
				},
			},
		},
	}
	
	common.DebugMessage("Creating network")
	netresp, err := cli.NetworkCreate(ctx, networkName, opts)
	
	if err != nil {
		return "", err
	}
	
	common.DebugMessage("ID: " + netresp.ID)
	if netresp.Warning != "" {
		common.DebugMessage("Warning: " + netresp.Warning)
	}
	
	return netresp.ID, err
}

// SafeCreateNetwork provides a wrapper to CreateNetwork, but first ensures that the
// network does not already exist.
func SafeCreateNetwork(cli *client.Client, networkName, subnet, iprange, gateway string, force, useExisting bool) (string, error) {
	
	// First we get a list of networks
	nets, err := GetNetworks(cli, false)
	if err != nil {
		return "", err
	}
	
	if netID, ok := nets[networkName]; ok {
		// Network already exists
		if useExisting {
			return netID, nil
		}
		if !force {
			return "", errors.New("the network " + networkName + " already exists")
		}
		
		common.DebugMessage("remove the existing network " + networkName)
		if err := RemoveNetwork(cli, netID); err != nil {
			return "", err
		}
	}
	return CreateNetwork(cli, networkName, subnet, iprange, gateway)
}

// RemoveNetwork removes a network
func RemoveNetwork(cli *client.Client, networkID string) error {
	ctx := context.Background()
	return cli.NetworkRemove(ctx, networkID)
}