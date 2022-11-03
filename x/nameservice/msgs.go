package nameservice

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetWhois struct {
	Name  string         //域名
	Value string         //域名解析值
	Owner sdk.AccAddress //域名所有者
}

func NewMsgSetWhois(name string, value string, owner sdk.AccAddress) MsgSetWhois {
	return MsgSetWhois{Name: name, Value: value, Owner: owner}
}

// 返回模块名称 并路由消息
func (msg MsgSetWhois) Route() string {
	return "MsgSetWhois nameservice"
}

// 返回消息的可读字符串/动作
func (msg MsgSetWhois) Type() string {
	return "MsgSetWhois Set whois"
}

// 不依赖链上状态对交易做基本的有效性验证，比如Fee的值是否合法有效，签名者数量和签名数量是否一致
func (msg MsgSetWhois) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("name or value can`t be empty")
	}
	return nil
}

// 将需要签名的消息进行有序的编码
func (msg MsgSetWhois) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// 获取可以签名的地址
func (msg MsgSetWhois) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgBuyWhois struct {
	Name  string         //所购买域名的名称
	Bid   sdk.Coins      //话费
	Buyer sdk.AccAddress //购买人
}

func NewMsgBuyWhois(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyWhois {
	return MsgBuyWhois{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}
func (msg MsgBuyWhois) Route() string {
	return "nameservice MsgBuyWhois"
}

func (msg MsgBuyWhois) Type() string {
	return "Buy whois"
}

func (msg MsgBuyWhois) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name connot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

func (msg MsgBuyWhois) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgBuyWhois) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
