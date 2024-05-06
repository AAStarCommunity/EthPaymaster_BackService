package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

var PaymasterAPIMethods = map[string]MethodFunctionFunc{}

type MethodFunctionFunc = func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error)

func init() {
	PaymasterAPIMethods["pm_sponsorUserOperation"] = TryPayUserOperationMethod()
	PaymasterAPIMethods["pm_supportEntrypoint"] = GetSupportEntryPointFunc()
	PaymasterAPIMethods["pm_strategyInfo"] = GetSupportStrategyFunc()
	PaymasterAPIMethods["pm_estimateUserOperationGas"] = EstimateUserOpGasFunc()
}

const (
	defaultStackSize = 4096
)

func getCurrentGoroutineStack() string {
	var buf [defaultStackSize]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

// Paymaster
// @Tags Paymaster
// @Description Paymaster JSON-RPC API
// @Accept json
// @Product json
// @Param rpcRequest body model.JsonRpcRequest true "JsonRpcRequest Model"
// @Router /api/v1/paymaster  [post]
// @Success 200
// @Security JWT
func Paymaster(ctx *gin.Context) {

	jsonRpcRequest := model.JsonRpcRequest{}
	response := model.GetResponse()

	defer func() {
		if r := recover(); r != nil {
			errInfo := fmt.Sprintf("[panic]: err : [%v] , stack :[%v]", r, getCurrentGoroutineStack())
			logrus.Error(errInfo)
			response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, fmt.Sprintf("%v", r))
		}

	}()

	if err := ctx.ShouldBindJSON(&jsonRpcRequest); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
	method := jsonRpcRequest.Method
	if method == "" {
		errStr := fmt.Sprintf("Request Error [method is empty]")
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}

	if methodFunc, ok := PaymasterAPIMethods[method]; ok {
		logrus.Debug(fmt.Sprintf("method: %s", method))
		result, err := methodFunc(ctx, jsonRpcRequest)
		logrus.Debugf("result: %v", result)
		if err != nil {
			errStr := fmt.Sprintf("%v", err)
			response.SetHttpCode(http.StatusInternalServerError).FailCode(ctx, http.StatusInternalServerError, errStr)
			return
		}
		response.WithData(result).Success(ctx)
		return
	} else {
		errStr := fmt.Sprintf("Request Error [method not found]")
		response.SetHttpCode(http.StatusBadRequest).FailCode(ctx, http.StatusBadRequest, errStr)
		return
	}
}
