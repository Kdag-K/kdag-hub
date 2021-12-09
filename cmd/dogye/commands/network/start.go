package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	config "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/docker"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/pelletier/go-toml"
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
	network := args[0]
	
	if err := startDockerNetwork(network); err != nil {
		return err
	}
	
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
	
	// Create a Docker Client.
	common.DebugMessage("Connecting to Docker Client")
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	
	// Create a Docker Network.
	common.DebugMessage(fmt.Sprintf("Created Network %s (%s)", conf.Docker.Name, networkID))
	
	// Next we build the docker configurations to get all of the configs ready to push.
	
	return nil
}

func exportDockerConfigs(conf *Config) error {
	
	// Configure some paths.
	netDir := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, conf.Network.Name)
	dockerDir := filepath.Join(netDir, dogyeDockerDir)
	err := files.CreateDirsIfNotExists([]string{dockerDir})
	if err != nil {
		return err
	}
	
	// loop around nodes
	for _, n := range conf.Nodes {
		if !n.NonNode {
			if err := exportDockerNodeConfig(netDir, dockerDir, &n); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func exportDockerNodeConfig(netDir, dockerDir string, node *node) error {
	netaddr := node.NetAddr
	if !strings.Contains(netaddr, ":") {
		netaddr += ":" + config.DefaultGossipPort
	}
	// Build output files.
	if node.Moniker != "" { // Should not be blank here, but safety first
		
		knodeDir := filepath.Join(dockerDir, node.Moniker, config.KnodeTomlDirDot)
		
		configDir := filepath.Join(knodeDir, config.ConfigDir)
		kdagConfigDir := filepath.Join(configDir, config.KdagDir)
		ethConfigDir := filepath.Join(configDir, config.EthDir)
		keystoreDir := filepath.Join(knodeDir, config.KeyStoreDir)
		
		common.DebugMessage("Creating config in " + knodeDir)
		
		err := files.CreateDirsIfNotExists([]string{
			kdagConfigDir,
			ethConfigDir,
			keystoreDir,
		})
		if err != nil {
			return err
		}
		
		copying := []copyRecord{
			{ // knode.toml
				from: filepath.Join(netDir, config.KnodeTomlFile),
				to:   filepath.Join(configDir, config.KnodeTomlFile),
			},
			{ // eth/genesis.json
				from: filepath.Join(netDir, config.GenesisJSON),
				to:   filepath.Join(ethConfigDir, config.GenesisJSON),
			},
			{ // kdag/peers.json
				from: filepath.Join(netDir, config.PeersJSON),
				to:   filepath.Join(kdagConfigDir, config.PeersJSON),
			},
			{ // kdag/peers.genesis.json
				from: filepath.Join(netDir, config.PeersGenesisJSON),
				to:   filepath.Join(kdagConfigDir, config.PeersGenesisJSON),
			},
			{ // keystore/<moniker>.json (private key)
				from: filepath.Join(netDir, config.KeyStoreDir, node.Moniker+".json"),
				to:   filepath.Join(keystoreDir, node.Moniker+".json"),
			},
			{ // keystore/<moniker>.text (password)
				from: filepath.Join(netDir, config.KeyStoreDir, node.Moniker+".txt"),
				to:   filepath.Join(keystoreDir, node.Moniker+".txt"),
			},
		}
		
		for _, f := range copying {
			files.CopyFileContents(f.from, f.to)
		}
		
		// Write a node description file containing all of the parameters needed
		// to start a container. Saves having to load and parse network.toml for
		//  every node
		nodeConfigFile := filepath.Join(dockerDir, node.Moniker+".toml")
		nodeConfig := dockerNodeConfig{
			Moniker: node.Moniker,
			NetAddr: strings.Split(netaddr, ":")[0],
		}
		
		tomlBytes, err := toml.Marshal(nodeConfig)
		if err != nil {
			return err
		}
		
		err = files.WriteToFile(nodeConfigFile, string(tomlBytes), files.OverwriteSilently)
		if err != nil {
			return err
		}
		
		// edit knode.toml and set kdag.listen appropriately
		
		// decrypt the validator private key, and dump it into the kdag config
		// dir (priv_key)
		
	}
	return nil
}