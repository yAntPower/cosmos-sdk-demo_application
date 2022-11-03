package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/auth"
	dbm "github.com/tendermint/tendermint/libs/db"
	tlog "github.com/tendermint/tendermint/libs/log"
)

const appName = "nameservice"

type nameServiceApp struct {
	*bam.BaseApp
}

func NewNameServiceApp(log tlog.Logger, db dbm.DB) *nameServiceApp {
	cdc := MakeCodec()
	bapp := NewBaseApp(appName, log, db, auth.DefaultTxDecoder(cdc))
	var app = &nameServiceApp{
		BaseApp: bapp,
		cdc:     cdc,
	}
	return app
}
