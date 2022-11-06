package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"nameservice/x/nameservice"
	"net/http"
)

const restName = "name"

// 提供程序访问模块的功能,定义REST客户端接口
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoisIsHandle(cdc, cliCtx, storeName)).Methods("GET")
}

// 获取解析
func resolveNameHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		paramType := vars[restName]
		res, err := cliCtx.QueryWithData(fmt.Sprintf("coustom/%s/resolve/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}

// 获取域名
func whoisIsHandle(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		paramsType := vars[restName]
		res, err := cliCtx.QueryWithData(fmt.Sprintf("coustom/%s/whois/%s", storeName, paramsType), nil)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}

func namesHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("coustom/%s/names", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}

type buyNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Amount  string       `json:"amount"`
	Buyer   string       `json:"buyer"`
}

// 购买处理
func buyNameHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req buyNameReq
		if !rest.ReadRESTReq(writer, request, cdc, &req) {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(writer) {
			return
		}
		addr, err := sdk.AccAddressFromBech32(req.Buyer)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		msg := nameservice.NewMsgBuyWhois(req.Name, coins, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.WriteGenerateStdTxResponse(writer, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type setNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Value   string       `json:"value"`
	Owner   string       `json:"owner"`
}

// 设置处理
func setNameHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req setNameReq
		if !rest.ReadRESTReq(writer, request, cdc, &req) {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(writer) {
			return
		}
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		msg := nameservice.NewMsgSetWhois(req.Name, req.Value, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.WriteGenerateStdTxResponse(writer, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
