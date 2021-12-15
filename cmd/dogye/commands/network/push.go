package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/docker"
	"github.com/pelletier/go-toml"
	
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
)

// Parameters for docker client
const (
	imgName     = "mosaicnetworks/monetd:latest"
	imgIsRemote = false
)

func newPushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [network] [node]",
		Short: "push a node onto the named network",
		Long: `
dogye network push

This command is called after 'dogye network start'. It builds a node based on
the configation files found for <node>, attaches it to the docker network, and
starts monetd.
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkPush,
	}
	
	return cmd
}

func networkPush(cmd *cobra.Command, args []string) error {
	networkName := args[0]
	nodeName := args[1]
	return pushDockerNode(networkName, nodeName, imgName, imgIsRemote)
}

// PushDockerNode builds a docker node, configures it and starts it.
func pushDockerNode(networkName, nodeName, imgName string, isRemoteImage bool) error {
	common.DebugMessage("Pushing network " + networkName + " node " + nodeName)
	
	// First we validate that the requested node has been created
	dockerpath := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, dogyeDockerDir)
	if !files.CheckIfExists(dockerpath) {
		return errors.New(" cannot find docker config for network " + networkName + ". Have you run dogye network start? ")
	}
	
	dockerconfigpath := filepath.Join(dockerpath, nodeName)
	if !files.CheckIfExists(dockerconfigpath) {
		return errors.New(" cannot find docker config folder for node " + nodeName)
	}
	
	dockerconfig := filepath.Join(dockerpath, nodeName+".toml")
	if !files.CheckIfExists(dockerconfig) {
		return errors.New(" cannot find docker config toml for node " + nodeName)
	}
	
	// Get Node Details
	// TODO:
	
	// Read key from file.
	tomlfile, err := ioutil.ReadFile(dockerconfig)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", dockerconfig, err)
	}
	
	config := dockerNodeConfig{}
	toml.Unmarshal(tomlfile, &config)
	
	common.DebugMessage("Container IP is " + config.NetAddr)
	
	// Start Docker Client
	common.DebugMessage("Connecting to Docker Client\n ")
	
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	
	// retrieve network ID based on network name
	var networkID string
	if nets, err := docker.GetNetworks(cli, false); err == nil {
		if net, ok := nets[networkName]; ok {
			networkID = net
		} else {
			return errors.New("network " + networkName + " is not running")
		}
	} else {
		common.ErrorMessage("Error getting network status")
		return nil
	}
	
	// Check current containers to see if node already exists
	containers, err := docker.GetContainers(cli, false)
	
	if existingNode, ok := containers[nodeName]; ok {
		return errors.New("node " + nodeName + " already exists (" + existingNode + ")")
	}
	
	// Create Node
	common.DebugMessage("Creating Container ")
	
	common.DebugMessage("Created Container " + containerID)
	
	// Copy Configuration to Node
	common.DebugMessage("Copying Config to Container ")
	
	// Configure Networking
	common.DebugMessage("Connecting Container to Network")
	
	// Start Node
	common.DebugMessage("Starting Container ")
	
	common.DebugMessage("Container Started")
	
	return nil
}