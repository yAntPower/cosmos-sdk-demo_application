package nameservice

import "github.com/cosmos/cosmos-sdk/codec"

// 注册本模块需要用Amino处理的类型
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBuyWhois{}, "nameservice/BuyWhois", nil)
	cdc.RegisterConcrete(MsgSetWhois{}, "nameservice/SetWhois", nil)
}
