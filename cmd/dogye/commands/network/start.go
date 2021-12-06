package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	
	"github.com/BurntSushi/toml"
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI flags.
var forceNetwork = false
var useExisting = false
var startNodes = false

type copyRecord struct {
	from string
	to   string
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
dogye network start

Starts a network. Does not start individual nodes.
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}
	
	addStartFlags(cmd)
	
	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&forceNetwork, "force-network", forceNetwork, "force network down if already exists")
	cmd.Flags().BoolVar(&useExisting, "use-existing", useExisting, "use existing network if already exists")
	cmd.Flags().BoolVar(&startNodes, "start-nodes", startNodes, "start nodes")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	return nil
}

// add startDockerNetwork to integrate this command with docker module.
func startDockerNetwork(networkName string) error {
	// Set some paths.
	thisNetworkDir := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)
	
	// Check expected config exists
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}
	
	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}
	
	var conf = Config{}
	
	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}
	
	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}
	
	common.DebugMessage("Configuring Network ", conf.Docker.Name)
	
	if conf.Docker.Name == "" {
		return errors.New("network " + networkName + " is not configured as a docker network")
	}
	
	// Create a Docker Client
	common.DebugMessage("Connecting to Docker Client")
	
	// Create a Docker Network.
	common.DebugMessage(fmt.Sprintf("Created Network %s (%s)", conf.Docker.Name, networkID))
	
	// Next we build the docker configurations to get all of the configs ready to push.
	
	return nil
}