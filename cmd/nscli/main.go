package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	app "nameservice"
	nsclient "nameservice/x/nameservice/client"
	nsrest "nameservice/x/nameservice/client/rest"
	"os"
	"path"
)

const (
	storeAcc = "acc"
	storeNS  = "nameservice"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.nscli")

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()
	mc := []sdk.ModuleClients{nsclient.NewModuleClient(storeNS, cdc)}
	rootCmd := &cobra.Command{
		Use:   "nscli",
		Short: "nameservice client",
	}
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return initConfig(rootCmd)
	}
	rootCmd.AddCommand(rpc.StatusCommand(), client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc, mc), txCmd(cdc, mc), client.LineBreak, lcd.ServeCommand(cdc, registerRoutes), client.LineBreak,
		keys.Commands(), client.LineBreak)
	executor := cli.PrepareMainCmd(rootCmd, "NS", defaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	rs.CliCtx = rs.CliCtx.WithAccountDecoder(rs.Cdc)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	auth.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeAcc)
	bank.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	nsrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeNS)
}
func queryCmd(cdc *codec.Codec, mc []sdk.ModuleClients) *cobra.Command {
	qCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying SubCommands",
	}
	qCmd.AddCommand(rpc.ValidatorCommand(cdc), rpc.BlockCommand(), tx.SearchTxCmd(cdc), tx.QueryTxCmd(cdc), client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc))
	for _, m := range mc {
		qCmd.AddCommand(m.GetQueryCmd())
	}
	return qCmd
}

func txCmd(cdc *codec.Codec, mc []sdk.ModuleClients) *cobra.Command {
	tCmd := &cobra.Command{
		Use:   "tx",
		Short: "transaction subcommands",
	}
	tCmd.AddCommand(bankcmd.SendTxCmd(cdc), client.LineBreak, authcmd.GetSignCommand(cdc), tx.GetBroadcastCommand(cdc),
		client.LineBreak)
	for _, m := range mc {
		tCmd.AddCommand(m.GetTxCmd())
	}
	return tCmd
}
func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}
	cfgFile := path.Join(home, "config", "config.toml")
	_, err = os.Stat(cfgFile)
	if err == nil {
		viper.SetConfigFile(cfgFile)
		if err = viper.ReadInConfig(); err != nil {
			return err
		}
	}
	err = viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID))
	if err != nil {
		return err
	}
	err = viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag))
	if err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
