package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	nameSerciveCli "nameservice/x/nameservice/client/cli"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// 返回支持的查询命令
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	nameSerQueryCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Querying commands for the nameService module",
	}
	nameSerQueryCmd.AddCommand(client.GetCommands(
		nameSerciveCli.GetCmdResolveName(mc.storeKey, mc.cdc),
		nameSerciveCli.GetCmdWhois(mc.storeKey, mc.cdc),
		nameSerciveCli.GetCmdNames(mc.storeKey, mc.cdc),
	)...)
	return nameSerQueryCmd
}

// 返回操作命令
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	nameSvcTxCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Nameservice transactions subcommands",
	}
	nameSvcTxCmd.AddCommand(client.GetCommands(nameSerciveCli.GetCmdBuyWhois(mc.cdc), nameSerciveCli.GetCmdSetWhois(mc.cdc))...)
	return nameSvcTxCmd
}
