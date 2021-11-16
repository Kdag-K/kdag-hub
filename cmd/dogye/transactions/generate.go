package transactions

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	
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
	
	return nil
}