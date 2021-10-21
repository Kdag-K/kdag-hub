package configuration

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/files"
)

const configTOML = `
# Set to true for extended logging
verbose = "{{.Verbose}}"

# The IP:PORT of the HTTP API service (defaults to :8080)
api-listen = "{{.APIAddr}}"

[kdag]

  # IP:PORT on the local machine where Kdag will bind its internal gossip
  # sockets. If this is not reachable from the outside, use 'advertise' to
  # define a routable address that other peers can reach.
  listen = "{{.Kdag.BindAddr}}"

  # IP:PORT advertised to other nodes. This is the address that other nodes use
  # to contact this node. It defaults to the listen address.
  # advertise = "{{.Kdag.AdvertiseAddr}}"

  # The heartbeat timer, the time in milliseconds between gossips
  heartbeat = "{{.Kdag.Heartbeat}}"

  # TCP timeout
  timeout = "{{.Kdag.TCPTimeout}}"

  # Number of items in the LRU cache
  cache-size = {{.Kdag.CacheSize}}

  # Max number of events per sync
  sync-limit = {{.Kdag.SyncLimit}}

  # Max number of pool connections
  max-pool = {{.Kdag.MaxPool}}

  # Bootstrap Kdag from database
  bootstrap = {{.Kdag.Bootstrap}}

  # Moniker for this node
  moniker = "{{.Kdag.Moniker}}"

  # Set to true to enable Maintenance Mode to start Kdag in a non-gossipping
  # suspended state.
  maintenance-mode = "{{.Kdag.MaintenanceMode}}"

[eth]
  # megabytes of memory allocated to internal caching
  # (min 16MB / database forced) (default 128)
  cache = {{.Eth.Cache}}
  
  # minimum gasprice of transactions submitted through this node (eg 1T) (default "0")
  min-gas-price = {{.Eth.MinGasPrice}}
`

// GlobalTOML converts the global Config object into a TOML string
func GlobalTOML() (string, error) {
	configTmpl, err := template.New("monetd.toml").Parse(configTOML)
	if err != nil {
		return "", fmt.Errorf("Error parsing monetd.toml template: %v", err)
	}
	
	var buf bytes.Buffer
	err = configTmpl.Execute(&buf, Global)
	if err != nil {
		return "", fmt.Errorf("Error executing monetd.toml template: %v", err)
	}
	
	return buf.String(), nil
}

// DumpGlobalTOML takes the global Config object, encodes it into a TOML string,
// and writes it to a file.
func DumpGlobalTOML(configDir, fileName string, interactive bool) error {
	tomlString, err := GlobalTOML()
	if err != nil {
		return err
	}
	
	options := files.OverwriteSilently
	
	if interactive {
		options = files.PromptIfExisting
	}
	
	if err := files.WriteToFile(
		filepath.Join(configDir, fileName),
		tomlString,
		options); err != nil {
		return err
	}
	
	return nil
}