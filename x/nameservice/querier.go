package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path, req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path, req, keeper)
		case QueryNames:
			return queryNames(ctx, path, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

type QueryResResolve struct {
	Value string `json:"value"`
}

func (q QueryResResolve) String() string {
	return q.Value
}
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]
	value := keeper.GetResolvesValue(ctx, name)
	if value == "" {
		errMsg := "could not resolve name." + name
		return []byte{}, sdk.ErrUnknownRequest(errMsg)
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
	if err != nil {
		panic("could not marshal result to JSON(queryResolve)")
	}
	return res, nil
}
func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]
	whois := keeper.GetWhois(ctx, name)
	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err != nil {
		panic("could not marshal result to JSON(queryWhois)")
	}
	return res, nil

}

type QueryResNames []string

func (qrn QueryResNames) String() string {
	return strings.Join(qrn[:], "\n")
}

func queryNames(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var names QueryResNames
	nameIter := keeper.GetNamesIterator(ctx)
	for ; nameIter.Valid(); nameIter.Next() {
		names = append(names, string(nameIter.Key()))
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, names)
	if err != nil {
		panic("could not marshal result to JSON(queryNames)")
	}
	return res, nil
}
