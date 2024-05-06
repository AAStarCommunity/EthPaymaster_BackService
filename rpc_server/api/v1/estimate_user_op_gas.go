package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/operator"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
)

// EstimateUserOpGas
// @Tags Sponsor
// @Description get the estimate gas of the userOp
// @Accept json
// @Product json
// @Router /api/v1/estimate-user-operation-gas [POST]
// @Param tryPay body model.UserOpRequest true "UserOp Request"
// @Success 200
// @Security JWT
func EstimateUserOpGas(c *gin.Context) {
	request := model.UserOpRequest{}
	response := model.GetResponse()
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
	if result, err := operator.GetEstimateUserOpGas(&request); err != nil {
		errStr := fmt.Sprintf("GetEstimateUserOpGas ERROR [%v]", err)
		response.SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
		return
	} else {
		response.WithDataSuccess(c, result)
		return
	}
}
func EstimateUserOpGasFunc() MethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest model.JsonRpcRequest) (result interface{}, err error) {
		request, err := ParseTryPayUserOperationParams(jsonRpcRequest.Params)
		if err != nil {
			return nil, xerrors.Errorf("ParseTryPayUserOperationParams ERROR [%v]", err)
		}
		if err := ValidateUserOpRequest(request); err != nil {
			return nil, xerrors.Errorf("Request Error [%v]", err)
		}
		if result, err := operator.GetEstimateUserOpGas(request); err != nil {
			return nil, xerrors.Errorf("GetEstimateUserOpGas ERROR [%v]", err)
		} else {
			return result, nil
		}
	}
}
