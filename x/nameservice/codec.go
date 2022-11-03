package nameservice

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBuyWhois{}, "nameservice/BuyWhois", nil)
	cdc.RegisterConcrete(MsgSetWhois{}, "nameservice/SetWhois", nil)
}
