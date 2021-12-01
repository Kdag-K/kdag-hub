package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/pelletier/go-toml"

	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"

	"github.com/spf13/cobra"
)

func newDumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump [network_name]",
		Short: "Dump the network settings",
		Long: `
dogye network dump
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkDump,
	}

	return cmd
}

func networkDump(cmd *cobra.Command, args []string) error {

	networkName = strings.TrimSpace(args[0])

	// Sanity check the network
	if !common.CheckMoniker(networkName) {
		return errors.New("the network name, " + networkName + ", is invalid")
	}

	if !files.CheckIfExists(filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName)) {
		return errors.New("the network, " + networkName + " has not been created")
	}

	networkTomlFile := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, networkTomlFileName)

	var conf = Config{}

	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}

	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}

	var dumpOut []string

	for _, n := range conf.Nodes {
		netaddr := n.NetAddr
		if idx := strings.Index(netaddr, ":"); idx > -1 {
			netaddr = netaddr[:idx]
		}

		dumpOut = append(dumpOut, n.Moniker+"|"+netaddr+"|"+n.Address+"|"+strconv.FormatBool(n.Validator)+"|"+strconv.FormatBool(n.NonNode))

	}

	for _, o := range dumpOut {
		fmt.Println(o)
	}

	return nil
}
