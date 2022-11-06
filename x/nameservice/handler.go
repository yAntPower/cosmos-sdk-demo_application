package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// 创建操作本模块的消息过滤器
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch m := msg.(type) {
		case MsgSetWhois:
			return handleMsgSetWhois(ctx, keeper, m)
		case MsgBuyWhois:
			return handleMsgBuyWhois(ctx, keeper, m)
		default:
			errmsg := fmt.Sprint("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errmsg).Result()

		}
	}
}
func handleMsgSetWhois(ctx sdk.Context, keeper Keeper, msg MsgSetWhois) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	keeper.SetResolvesValue(ctx, msg.Name, msg.Value)
	return sdk.Result{}
}

func handleMsgBuyWhois(ctx sdk.Context, keeper Keeper, msg MsgBuyWhois) sdk.Result {
	if keeper.GetPrice(ctx, msg.Name).IsAllGTE(msg.Bid) {
		return sdk.ErrUnauthorized("Bid not high enough").Result()
	}
	if keeper.HasOwner(ctx, msg.Name) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer have not enougs Coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) //没有接收者燃烧掉
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer have not enougs Coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}
