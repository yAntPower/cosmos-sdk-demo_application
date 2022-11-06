package app

import (
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"nameservice/x/nameservice"
)

const appName = "nameservice"

type nameServiceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyNS            *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tKeyParams       *sdk.TransientStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper        params.Keeper
	nsKeeper            nameservice.Keeper
}

func NewNameServiceApp(log tlog.Logger, db dbm.DB) *nameServiceApp {
	cdc := MakeCodec()
	bapp := bam.NewBaseApp(appName, log, db, auth.DefaultTxDecoder(cdc))
	var app = &nameServiceApp{
		BaseApp:          bapp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyNS:            sdk.NewKVStoreKey("ns"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tKeyParams:       sdk.NewTransientStoreKey("transient_params"),
	}
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tKeyParams)
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)
	app.nsKeeper = nameservice.NewKeeper(app.bankKeeper, app.keyNS, app.cdc)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
	app.Router().AddRoute("bank", bank.NewHandler(app.bankKeeper)).AddRoute("nameservice", nameservice.NewHandler(app.nsKeeper))
	app.QueryRouter().AddRoute("nameservice", nameservice.NewQuerier(app.nsKeeper)).AddRoute("acc", auth.NewQuerier(app.accountKeeper))
	app.SetInitChainer(app.initChainer)
	app.MountStores(app.keyMain, app.keyAccount, app.keyNS, app.keyFeeCollection, app.keyParams, app.tKeyParams)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

type GenesisState struct {
	AuthData auth.GenesisState   `json:"auth"`
	BankData bank.GenesisState   `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
}

// 初始化创世文件中的帐户
func (app *nameServiceApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJson := req.AppStateBytes
	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJson, genesisState)
	if err != nil {
		panic(err)
	}
	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	return abci.ResponseInitChain{}
}
func (app *nameServiceApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}
	appendAcountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}
	app.accountKeeper.IterateAccounts(ctx, appendAcountsFn)
	genState := GenesisState{Accounts: accounts, AuthData: auth.DefaultGenesisState(), BankData: bank.DefaultGenesisState()}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	return
}

// 把所有用到的模块内的数据类型都拿去注册codec，目的是为了可以使用Amino进行编解码
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	nameservice.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
