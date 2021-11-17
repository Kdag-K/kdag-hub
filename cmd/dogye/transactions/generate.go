package transactions

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/BurntSushi/toml"
    config "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/cmd/dogye/network"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate transactions",
		Args:  cobra.ArbitraryArgs,
		RunE:  generateTransactions,
	}
	
	addGenerateFlags(cmd)
	return cmd
}
func addGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&networkName, "network", "n", "", "network name")
	cmd.Flags().StringVar(&ips, "ips", "", "ips.dat file path")
	cmd.Flags().StringVar(&faucet, "faucet", faucet, "faucet account moniker")
	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to generate")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	viper.BindPFlags(cmd.Flags())
}

func loadIPS() (map[string]string, error) {
	rtn := make(ipmapping)
	
	if !files.CheckIfExists(ips) {
		return nil, errors.New("ip mapping does not exist: " + ips)
	}
	
	file, err := os.Open(ips)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		
		ippair := strings.Fields(scanner.Text())
		if len(ippair) != 2 {
			return nil, errors.New("malformed ip mapping " + strconv.Itoa(len(ippair)))
		}
		rtn[ippair[0]] = ippair[1]
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return rtn, nil
}

func generateTransactions(cmd *cobra.Command, args []string) error {
	var ipmap ipmapping
	
	var surplusCreditBig = new(big.Int).SetInt64(int64(surplusCredit))
	
	if !common.CheckMoniker(networkName) {
		return errors.New("network name must only contains characters in the range 0-9 or A-Z or a-z")
	}
	
	networkDir := filepath.Join(_dogye, dogyeNetworksDir, networkName)
	if !files.CheckIfExists(networkDir) {
		return errors.New("network does not exist: " + networkDir)
	}
	
	keystore := filepath.Join(networkDir, dogyeKeystoreDir)
	if !files.CheckIfExists(keystore) {
		return errors.New("keystore does not exist: " + keystore)
	}
	
	networkTomlFile := filepath.Join(networkDir, networkTomlFileName)
	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("toml file does not exist: " + networkTomlFile)
	}
	
	transDir := filepath.Join(networkDir, dogyeTransactionsDir)
	if err := files.SafeRename(transDir); err != nil {
		return err
	}
	
	if err := files.CreateDirsIfNotExists([]string{transDir}); err != nil {
		return err
	}
	faucetTransFile := filepath.Join(transDir, "faucet.json")
	transFile := filepath.Join(transDir, "trans.json")
	deltaFile := filepath.Join(transDir, "delta.json")
	
	var conf = network.Config{}
	
	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}
	
	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}
	
	if ips != "" {
		ipmap, err = loadIPS()
		if err != nil {
			return err
		}
	}
	
	var nodes []node
	var accounts []account
	var debits []*big.Int
	var credits []*big.Int
	var faucetAccount *account
	var trans []transaction
	var deltas []delta
	var nodeTrans []nodeTransactions
	common.DebugMessage("Parsing network.toml for node and accounts")
	for _, n := range conf.Nodes {
		
		netaddr := n.NetAddr
		moniker := n.Moniker
		balance := n.Tokens
		
		if netaddr == "" {
			if moniker == faucet {
				faucetAccount = &account{
					Address:   n.Address,
					Moniker:   n.Moniker,
					PubKeyHex: n.PubKeyHex,
					Tokens:    balance,
				}
				
				common.DebugMessage("faucet ", faucetAccount.Moniker, balance)
			} else {
				accounts = append(accounts, account{
					Address:   n.Address,
					Moniker:   n.Moniker,
					PubKeyHex: n.PubKeyHex,
					Tokens:    balance,
				})
				
				credits = append(credits, new(big.Int))
				debits = append(debits, new(big.Int))
				
				common.DebugMessage("account ", moniker, balance)
			}
		} else {
			if ipmap != nil {
				netaddr = ipmap[netaddr]
			}
			if !strings.Contains(netaddr, ":") && len(netaddr) > 0 {
				netaddr += config.DefaultAPIAddr
			}
			nodes = append(nodes, node{
				NetAddr: netaddr,
				Moniker: moniker})
			common.DebugMessage("node ", moniker)
		}
		
	}
	
	nodecnt := len(nodes)
	accountcnt := len(accounts)
	
	common.InfoMessage("Read " + strconv.Itoa(nodecnt) + " nodes.")
	common.InfoMessage("Read " + strconv.Itoa(accountcnt) + " accounts.")
	
	if faucetAccount == nil {
		return errors.New("faucet account not found: " + faucet)
	}
	
	if accountcnt < 2 {
		return errors.New("you must have at least 2 accounts")
	}
	
	common.InfoMessage("Faucet account, " + faucet + ", found.")
	
	for i := 0; i < accountcnt; i++ {
		nodeTrans = append(nodeTrans, nodeTransactions{
			Address: accounts[i].Address,
			Moniker: accounts[i].Moniker,
		})
	}
	
	
	return nil
}