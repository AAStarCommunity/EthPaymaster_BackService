package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"net/http"
)

var PaymasterAPIMethods = map[string]MethodFunctionFunc{}

type MethodFunctionFunc = func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error)

func init() {
	PaymasterAPIMethods["pm_sponsorUserOperation"] = TryPayUserOperationMethod()
	PaymasterAPIMethods["pm_supportEntrypoint"] = GetSupportEntryPointFunc()
	PaymasterAPIMethods["pm_estimateUserOperationGas"] = EstimateUserOpGasFunc()
	PaymasterAPIMethods["pm_paymasterAccount"] = GetSupportPaymaster()
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
			errInfo := fmt.Sprintf("[panic]: err : [%v] , stack :[%v]", r, utils.GetCurrentGoroutineStack())
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

func GetSupportPaymaster() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		if jsonRpcRequest.Params[0] == nil {
			return nil, xerrors.Errorf("Request Error [network is empty]")
		}
		networkStr, ok := jsonRpcRequest.Params[0].(string)
		if !ok {
			return nil, xerrors.Errorf("Request Error [network is not string]")
		}
		paymasterSet, err := config.GetSupportPaymaster(global_const.Network(networkStr))
		if err != nil {
			return nil, err
		}
		return paymasterSet.ToSlice(), nil
	}
}

func GetSupportEntryPointFunc() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		if jsonRpcRequest.Params[0] == nil {
			return nil, xerrors.Errorf("Request Error [network is empty]")
		}
		networkStr, ok := jsonRpcRequest.Params[0].(string)
		if !ok {
			return nil, xerrors.Errorf("Request Error [network is not string]")
		}
		entryPoints, err := config.GetSupportEntryPoints(global_const.Network(networkStr))
		if err != nil {
			return nil, err
		}
		return entryPoints.ToSlice(), nil
	}
}
func EstimateUserOpGasFunc() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		request, err := parseTryPayUserOperationParams(jsonRpcRequest.Params)
		if err != nil {
			return nil, xerrors.Errorf("parseTryPayUserOperationParams ERROR [%v]", err)
		}
		if err := validateUserOpRequest(request); err != nil {
			return nil, xerrors.Errorf("Request Error [%v]", err)
		}
		if result, err := operator.GetEstimateUserOpGas(request); err != nil {
			return nil, xerrors.Errorf("GetEstimateUserOpGas ERROR [%v]", err)
		} else {
			return result, nil
		}
	}
}

func TryPayUserOperationMethod() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		request, err := parseTryPayUserOperationParams(jsonRpcRequest.Params)
		logrus.Debug("parseTryPayUserOperationParams result: ", request)

		if err != nil {
			return nil, xerrors.Errorf("parseTryPayUserOperationParams ERROR [%v]", err)
		}
		if err := validateUserOpRequest(request); err != nil {
			return nil, xerrors.Errorf("Request Error [%v]", err)
		}
		logrus.Debugf("After Validate ")

		if result, err := operator.TryPayUserOpExecute(request); err != nil {
			return nil, xerrors.Errorf("TryPayUserOpExecute ERROR [%v]", err)
		} else {
			return result, nil
		}
	}
}
func parseTryPayUserOperationParams(params []interface{}) (*model.UserOpRequest, error) {
	if len(params) < 2 {
		return nil, xerrors.Errorf("params length is less than 2")
	}
	result := model.UserOpRequest{}
	userInputParam := params[0]
	if userInputParam == nil {
		return nil, xerrors.Errorf("user input is nil")
	}
	userOpInput := userInputParam.(map[string]any)
	result.UserOp = userOpInput

	extraParam := params[1]
	if extraParam == nil {
		return nil, xerrors.Errorf("extra is nil")
	}
	extra := extraParam.(map[string]any)
	if extra["strategy_code"] != nil {
		result.ForceStrategyId = extra["strategy_code"].(string)
	}
	if extra["network"] != nil {
		result.ForceNetwork = extra["network"].(global_const.Network)
	}
	if extra["token"] != nil {
		result.Erc20Token = extra["token"].(global_const.TokenType)
	}
	if extra["version"] != nil {
		result.EntryPointVersion = extra["version"].(global_const.EntrypointVersion)
	}
	return &result, nil
}

func validateUserOpRequest(request *model.UserOpRequest) error {
	if request.ForceStrategyId != "" {
		return nil
	}
	if request.ForceNetwork == "" {
		return xerrors.Errorf("ForceNetwork is empty")
	}

	return nil
}