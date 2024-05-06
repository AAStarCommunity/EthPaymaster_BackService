package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"net/http"
)

// TryPayUserOperation
// @Tags Sponsor
// @Description sponsor the userOp
// @Accept json
// @Product json
// @Router /api/v1/try-pay-user-operation [POST]
// @Param tryPay body model.UserOpRequest true "UserOp Request"
// @Success 200
// @Security JWT
func TryPayUserOperation(c *gin.Context) {
	request := model.UserOpRequest{}
	response := model.GetResponse()

	//1. API validate
	if err := c.ShouldBindJSON(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	if err := ValidateUserOpRequest(&request); err != nil {
		errStr := fmt.Sprintf("Request Error [%v]", err)
		response.SetHttpCode(http.StatusBadRequest).FailCode(c, http.StatusBadRequest, errStr)
		return
	}

	//2. recall service
	if result, err := operator.TryPayUserOpExecute(&request); err != nil {
		errStr := fmt.Sprintf("TryPayUserOpExecute ERROR [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	} else {
		response.WithDataSuccess(c, result)
		return
	}
}

func TryPayUserOperationMethod() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		request, err := ParseTryPayUserOperationParams(jsonRpcRequest.Params)
		logrus.Debug("ParseTryPayUserOperationParams result: ", request)

		if err != nil {
			return nil, xerrors.Errorf("ParseTryPayUserOperationParams ERROR [%v]", err)
		}
		if err := ValidateUserOpRequest(request); err != nil {
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
func ParseTryPayUserOperationParams(params []interface{}) (*model.UserOpRequest, error) {
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

func ValidateUserOpRequest(request *model.UserOpRequest) error {
	if request.ForceStrategyId != "" {
		return nil
	}
	if request.ForceNetwork == "" {
		return xerrors.Errorf("ForceNetwork is empty")
	}

	return nil
}
