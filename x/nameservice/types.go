package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// 习惯将自有类型定义在types中
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("xingdao", 1)}

func NewWhois() Whois {
	return Whois{Price: MinNamePrice}
}
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner:%s
Value:%s
Price:%s`, w.Owner, w.Value, w.Price))
}
