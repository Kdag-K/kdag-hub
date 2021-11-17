package transactions

import (
	"bufio"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	
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
	
	return nil
}